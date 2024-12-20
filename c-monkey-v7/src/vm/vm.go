package vm

import (
	"Compiler/c-monkey-v7/src/code"
	"Compiler/c-monkey-v7/src/compiler"
	"Compiler/c-monkey-v7/src/object"
	"fmt"
)

const StackSize = 2048
const GlobalsSize = 65536
const MaxFrames = 1024

// Both these values are immutable so we can get away with create global versions instead of creating a new *object.Boolean each time
// This also makes comparisons easier without having to unpack pointers, ie, just compare object.Object == any of the following
var True = &object.Boolean{Value: true}
var False = &object.Boolean{Value: false}
var Null = &object.Null{}

type VM struct {
	constants []object.Object // Generated by compiler

	stack   []object.Object
	sp      int             // Always points to the next free value, not current value. Top of stack is stack[sp-1]
	globals []object.Object // Global identifier store

	frames      []*Frame // Frames for the execution of functions (Look at v6-Notes.md to understand what frames are)
	framesIndex int      // Like sp, points to the next free value, not the current value
}

func (vm *VM) currentFrame() *Frame {
	return vm.frames[vm.framesIndex-1]
}

func (vm *VM) pushFrame(f *Frame) {
	vm.frames[vm.framesIndex] = f
	vm.framesIndex++
}

func (vm *VM) popFrame() *Frame {
	vm.framesIndex--
	return vm.frames[vm.framesIndex]
}

func New(bytecode *compiler.Bytecode) *VM {
	mainFn := &object.CompiledFunction{Instructions: bytecode.Instructions} // Treating main() as a function on its own
	mainFrame := NewFrame(mainFn, 0)                                        // Creating a function for main

	frames := make([]*Frame, MaxFrames) // Creating a frame for the main
	frames[0] = mainFrame               // Main function is the first frame

	return &VM{
		constants: bytecode.Constants,

		stack: make([]object.Object, StackSize),
		sp:    0,

		globals: make([]object.Object, GlobalsSize),

		frames:      frames,
		framesIndex: 1, // Pointing to the next empty index, not the actual "top"
	}
}

// Constructor to maintain globals stores and compiler bytecode across executions
func NewWithGlobalStore(bytecode *compiler.Bytecode, s []object.Object) *VM {
	vm := New(bytecode)
	vm.globals = s
	return vm
}

func (vm *VM) LastPoppedStackElem() object.Object {
	return vm.stack[vm.sp]
}

func (vm *VM) Run() error {
	var ip int
	var ins code.Instructions
	var op code.Opcode

	for vm.currentFrame().ip < len(vm.currentFrame().Instructions())-1 {
		vm.currentFrame().ip++ // Increment per loop, control var

		ip = vm.currentFrame().ip              // Just to make the rest of the code easier to read, storing ip in local var
		ins = vm.currentFrame().Instructions() // Same reason as above. Easier for readability
		op = code.Opcode(ins[ip])              // For easier readability

		switch op {
		case code.OpConstant:
			constIndex := code.ReadUint16(ins[ip+1:]) // Decode step, starting with the byte right after the opcode
			vm.currentFrame().ip += 2                 // Incrementing by 2 because this is the width
			err := vm.push(vm.constants[constIndex])  // Execute step
			if err != nil {
				return err
			}

		case code.OpAdd, code.OpSub, code.OpDiv, code.OpMul:
			err := vm.executeBinaryOperation(op)
			if err != nil {
				return err
			}

		case code.OpPop:
			vm.pop()

		case code.OpTrue:
			err := vm.push(True)
			if err != nil {
				return err
			}

		case code.OpFalse:
			err := vm.push(False)
			if err != nil {
				return err
			}

		case code.OpEqual, code.OpNotEqual, code.OpGreaterThan, code.OpGreaterThanOrEqual:
			err := vm.executeComparison(op)
			if err != nil {
				return err
			}

		case code.OpBang:
			err := vm.executeBangOperator()
			if err != nil {
				return err
			}

		case code.OpMinus:
			err := vm.executeMinusOperator()
			if err != nil {
				return err
			}

		case code.OpJump:
			pos := int(code.ReadUint16(ins[ip+1:])) // Read the operand right next to the OpCode
			vm.currentFrame().ip = pos - 1          // set instruction pointer to the target of our jump. We do -1 so that the increment of the for loop can actually get us to the target

		case code.OpJumpNotTruthy:
			pos := int(code.ReadUint16(ins[ip+1:])) // Read operand right next to Opcode
			vm.currentFrame().ip += 2               // Increment to not read operand again in the next cycle

			condition := vm.pop()     // Pop the result of the condition
			if !isTruthy(condition) { // Jump if not truthy
				vm.currentFrame().ip = pos - 1
			}

		case code.OpNull:
			err := vm.push(Null)
			if err != nil {
				return err
			}

		case code.OpSetGlobal:
			globalIndex := code.ReadUint16(ins[ip+1:]) // Read operand
			vm.currentFrame().ip += 2                  // Increment to not read operand in next cycle

			vm.globals[globalIndex] = vm.pop() // Pop the value from the stack and assign it in the globals store

		case code.OpGetGlobal:
			globalIndex := code.ReadUint16(ins[ip+1:])
			vm.currentFrame().ip += 2

			err := vm.push(vm.globals[globalIndex])
			if err != nil {
				return err
			}

		case code.OpArray:
			numElements := int(code.ReadUint16(ins[ip+1:]))
			vm.currentFrame().ip += 2

			array := vm.buildArray(vm.sp-numElements, vm.sp)
			vm.sp = vm.sp - numElements

			err := vm.push(array)
			if err != nil {
				return err
			}

		case code.OpHash:
			numElements := int(code.ReadUint16(ins[ip+1:]))
			vm.currentFrame().ip += 2

			hash, err := vm.buildHash(vm.sp-numElements, vm.sp)
			if err != nil {
				return err
			}
			vm.sp = vm.sp - numElements

			err = vm.push(hash)
			if err != nil {
				return err
			}

		case code.OpIndex:
			index := vm.pop()
			left := vm.pop()

			err := vm.executeIndexExpression(left, index)
			if err != nil {
				return err
			}

		case code.OpCall:
			numArgs := code.ReadUint8(ins[ip+1:])
			vm.currentFrame().ip += 1

			err := vm.callFunction(int(numArgs))
			if err != nil {
				return err
			}

		case code.OpReturnValue:
			returnValue := vm.pop() // The latest element on the stack should be the return value

			frame := vm.popFrame()        // Pop the latest executed function, so that the next execution happens in the main flow
			vm.sp = frame.basePointer - 1 // Reset/clean the stack after function execution

			err := vm.push(returnValue)
			if err != nil {
				return err
			}

		case code.OpReturn:
			frame := vm.popFrame()
			vm.sp = frame.basePointer - 1

			err := vm.push(Null)
			if err != nil {
				return err
			}

		case code.OpSetLocal:
			localIndex := code.ReadUint8(ins[ip+1:]) // Read operand, ie, the index in this case
			vm.currentFrame().ip += 1                // Increment to not read operand in the next cycle

			frame := vm.currentFrame()
			vm.stack[frame.basePointer+int(localIndex)] = vm.pop() // Assign value in the stack with index as the offset

		case code.OpGetLocal:
			localIndex := code.ReadUint8(ins[ip+1:]) // Read operand
			vm.currentFrame().ip += 1                // Increment to not read operand again in the next cycle

			frame := vm.currentFrame()

			err := vm.push(vm.stack[frame.basePointer+int(localIndex)])
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (vm *VM) push(o object.Object) error {
	if vm.sp >= StackSize {
		return fmt.Errorf("stack overflow")
	}

	vm.stack[vm.sp] = o
	vm.sp++

	return nil
}

func (vm *VM) pop() object.Object {
	o := vm.stack[vm.sp-1]
	vm.sp--
	return o
}

func isTruthy(obj object.Object) bool {
	switch obj := obj.(type) {
	case *object.Boolean:
		return obj.Value
	case *object.Null:
		return false
	default:
		return true
	}
}

func (vm *VM) executeBinaryOperation(op code.Opcode) error {
	right := vm.pop()
	left := vm.pop()

	leftType := left.Type()
	rightType := right.Type()

	switch {
	case leftType == object.INTEGER_OBJ && rightType == object.INTEGER_OBJ:
		return vm.executeBinaryIntegerOperation(op, left, right)
	case leftType == object.STRING_OBJ && rightType == object.STRING_OBJ:
		return vm.executeBinaryStringOperation(op, left, right)
	default:
		return fmt.Errorf("unsupported types for binary operation: %s %s", leftType, rightType)
	}
}

func (vm *VM) executeBinaryIntegerOperation(op code.Opcode, left, right object.Object) error {
	leftValue := left.(*object.Integer).Value
	rightValue := right.(*object.Integer).Value

	var result int64

	switch op {
	case code.OpAdd:
		result = leftValue + rightValue
	case code.OpSub:
		result = leftValue - rightValue
	case code.OpMul:
		result = leftValue * rightValue
	case code.OpDiv:
		result = leftValue / rightValue
	default:
		return fmt.Errorf("unknown integer operator: %d", op)
	}

	return vm.push(&object.Integer{Value: result})
}

func (vm *VM) executeBinaryStringOperation(op code.Opcode, left, right object.Object) error {
	if op != code.OpAdd {
		return fmt.Errorf("unknown string operator: %d", op)
	}

	leftValue := left.(*object.String).Value
	rightValue := right.(*object.String).Value

	return vm.push(&object.String{Value: leftValue + rightValue})
}

func (vm *VM) executeComparison(op code.Opcode) error {
	right := vm.pop() // Right value always popped first, since it was pushed most recently, ie. LIFO
	left := vm.pop()

	if left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ {
		return vm.executeIntegerComparison(op, left, right)
	}

	switch op {
	case code.OpEqual:
		return vm.push(nativeBooltoBooleanObject(left == right))
	case code.OpNotEqual:
		return vm.push(nativeBooltoBooleanObject(left != right))
	default:
		return fmt.Errorf("unknown operator: %d (%s %s)", op, left.Type(), right.Type())
	}
}

func (vm *VM) executeIntegerComparison(op code.Opcode, left, right object.Object) error {
	rightValue := right.(*object.Integer).Value
	leftValue := left.(*object.Integer).Value

	switch op {
	case code.OpEqual:
		return vm.push(nativeBooltoBooleanObject(leftValue == rightValue))
	case code.OpNotEqual:
		return vm.push(nativeBooltoBooleanObject(leftValue != rightValue))
	case code.OpGreaterThan:
		return vm.push(nativeBooltoBooleanObject(leftValue > rightValue))
	case code.OpGreaterThanOrEqual:
		return vm.push(nativeBooltoBooleanObject(leftValue >= rightValue))
	default:
		return fmt.Errorf("unknown operator: %d", op)
	}
}

func nativeBooltoBooleanObject(input bool) *object.Boolean {
	if input {
		return True
	}

	return False
}

func (vm *VM) executeBangOperator() error {
	operand := vm.pop()

	switch operand {
	case True:
		return vm.push(False)
	case False:
		return vm.push(True)
	case Null:
		return vm.push(True)
	default:
		return vm.push(False)
	}
}

func (vm *VM) executeMinusOperator() error {
	operand := vm.pop()

	if operand.Type() != object.INTEGER_OBJ {
		return fmt.Errorf("unsupported type for negation: %s", operand.Type())
	}

	value := operand.(*object.Integer).Value
	return vm.push(&object.Integer{Value: -value})
}

func (vm *VM) executeIndexExpression(left, index object.Object) error {
	switch {
	case left.Type() == object.ARRAY_OBJ && index.Type() == object.INTEGER_OBJ:
		return vm.executeArrayIndex(left, index)
	case left.Type() == object.HASH_OBJ:
		return vm.executeHashIndex(left, index)
	default:
		return fmt.Errorf("index operator not supported: %s", left.Type())
	}
}

func (vm *VM) executeArrayIndex(array, index object.Object) error {
	arrayObject := array.(*object.Array)
	i := index.(*object.Integer).Value
	max := int64(len(arrayObject.Elements) - 1)

	if i < 0 || i > max {
		return vm.push(Null)
	}

	return vm.push(arrayObject.Elements[i])
}

func (vm *VM) executeHashIndex(hash, index object.Object) error {
	hashObject := hash.(*object.Hash)

	key, ok := index.(object.Hashable)
	if !ok {
		return fmt.Errorf("unusable as hash key: %s", index.Type())
	}

	pair, ok := hashObject.Pairs[key.HashKey()]
	if !ok {
		return vm.push(Null)
	}

	return vm.push(pair.Value)
}

func (vm *VM) buildArray(startIndex, endIndex int) object.Object {
	elements := make([]object.Object, endIndex-startIndex)

	for i := startIndex; i < endIndex; i++ {
		elements[i-startIndex] = vm.stack[i]
	}

	return &object.Array{Elements: elements}
}

func (vm *VM) buildHash(startIndex, endIndex int) (object.Object, error) {
	hashedPairs := make(map[object.HashKey]object.HashPair)

	for i := startIndex; i < endIndex; i += 2 {
		key := vm.stack[i]
		value := vm.stack[i+1]

		pair := object.HashPair{Key: key, Value: value}
		hashKey, ok := key.(object.Hashable)
		if !ok {
			return nil, fmt.Errorf("unusable as hash key: %s", key.Type())
		}

		hashedPairs[hashKey.HashKey()] = pair
	}

	return &object.Hash{Pairs: hashedPairs}, nil
}

func (vm *VM) callFunction(numArgs int) error {

	fn, ok := vm.stack[vm.sp-1-numArgs].(*object.CompiledFunction)
	if !ok {
		return fmt.Errorf("calling non-function")
	}

	if numArgs != fn.NumParameters {
		return fmt.Errorf("wrong number of arguments: want=%d, got=%d", fn.NumParameters, numArgs)
	}
	frame := NewFrame(fn, vm.sp-numArgs) // Get new frame from frame.go

	/* Push frame on to VM frame stack frame.
	Essentially, this modifies the currentFrame() that the VM is in, so after this happens,
	the rest of the execution will be function first
	When the frame is finally popped as part of the return statement executions, that is when the flow will return to the
	original frame of the execution */
	vm.pushFrame(frame)
	vm.sp = frame.basePointer + fn.NumLocals // Creating a "hole" for local bindings by incrementing the stack pointer NumLocal times

	return nil
}
