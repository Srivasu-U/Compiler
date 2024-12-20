package compiler

import (
	"Compiler/c-monkey-v3/src/ast"
	"Compiler/c-monkey-v3/src/code"
	"Compiler/c-monkey-v3/src/object"
	"fmt"
)

type Compiler struct {
	instructions        code.Instructions  // Generated bytecode
	constants           []object.Object    // Internal constant pool
	lastInstruction     EmittedInstruction // Last instruction
	previousInstruction EmittedInstruction // The one before last
}

type Bytecode struct { // Both are exportable fields since they start with capitalized letters. This gets passed into the VM
	Instructions code.Instructions
	Constants    []object.Object
}

type EmittedInstruction struct {
	Opcode   code.Opcode
	Position int // We can use the Position to check where to jump to in case of if...else, for example
}

func New() *Compiler {
	return &Compiler{
		instructions:        code.Instructions{},
		constants:           []object.Object{},
		lastInstruction:     EmittedInstruction{},
		previousInstruction: EmittedInstruction{},
	}
}

func (c *Compiler) Compile(node ast.Node) error {
	switch node := node.(type) {
	case *ast.Program:
		for _, s := range node.Statements {
			err := c.Compile(s)
			if err != nil {
				return err
			}
		}

	case *ast.ExpressionStatement:
		err := c.Compile(node.Expression)
		if err != nil {
			return err
		}
		c.emit(code.OpPop)

	case *ast.InfixExpression:
		// if < or <= reorder the execution of stack push for operands, ie, right first and then left
		if node.Operator == "<" || node.Operator == "<=" {
			err := c.Compile(node.Right)
			if err != nil {
				return err
			}

			err = c.Compile(node.Left)
			if err != nil {
				return err
			}

			if node.Operator == "<" {
				c.emit(code.OpGreaterThan)
			} else {
				c.emit(code.OpGreaterThanOrEqual)
			}
			return nil
		}
		err := c.Compile(node.Left)
		if err != nil {
			return err
		}

		err = c.Compile(node.Right)
		if err != nil {
			return err
		}

		switch node.Operator {
		case "+":
			c.emit(code.OpAdd)
		case "-":
			c.emit(code.OpSub)
		case "*":
			c.emit(code.OpMul)
		case "/":
			c.emit(code.OpDiv)
		case ">":
			c.emit(code.OpGreaterThan)
		case ">=":
			c.emit(code.OpGreaterThanOrEqual)
		case "==":
			c.emit(code.OpEqual)
		case "!=":
			c.emit(code.OpNotEqual)
		default:
			return fmt.Errorf("Unknown operator %s", node.Operator)
		}

	case *ast.IntegerLiteral:
		integer := &object.Integer{Value: node.Value}
		c.emit(code.OpConstant, c.addConstant(integer))

	case *ast.Boolean:
		if node.Value {
			c.emit(code.OpTrue)
		} else {
			c.emit(code.OpFalse)
		}

	case *ast.PrefixExpression:
		err := c.Compile(node.Right)
		if err != nil {
			return err
		}

		switch node.Operator {
		case "!":
			c.emit(code.OpBang)
		case "-":
			c.emit(code.OpMinus)
		default:
			return fmt.Errorf("unknown operator %s", node.Operator)
		}

	case *ast.IfExpression:
		err := c.Compile(node.Condition)
		if err != nil {
			return err
		}

		jumpNotTruthyPos := c.emit(code.OpJumpNotTruthy, 9999) // Jump location if the condition is false

		err = c.Compile(node.Consequence)
		if err != nil {
			return err
		}

		if c.lastInstructionIsPop() { // We remove the last pop because Block expressions always produce an OpPop which is redundant and we need to keep the final value that the block evaluates to
			c.removeLastPop()
		}

		// Jump location that is always executed to skip over the alternative
		jumpPos := c.emit(code.OpJump, 9999)

		// This basically gives us the offset after the execution of Consequence is done and where we want to jump to
		afterConsequencePos := len(c.instructions)
		// Going back and changing the operand, ie, the position where we jump to, instead of 9999
		c.changeOperand(jumpNotTruthyPos, afterConsequencePos)

		if node.Alternative == nil {
			c.emit(code.OpNull)
		} else {
			err := c.Compile(node.Alternative)
			if err != nil {
				return err
			}

			if c.lastInstructionIsPop() {
				c.removeLastPop()
			}
		}
		afterAlternativePos := len(c.instructions)
		c.changeOperand(jumpPos, afterAlternativePos)

	case *ast.BlockStatement:
		for _, s := range node.Statements {
			err := c.Compile(s)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (c *Compiler) Bytecode() *Bytecode {
	return &Bytecode{
		Instructions: c.instructions,
		Constants:    c.constants,
	}
}

// Adding constant to constant pool
func (c *Compiler) addConstant(obj object.Object) int {
	c.constants = append(c.constants, obj)
	return len(c.constants) - 1
}

func (c *Compiler) emit(op code.Opcode, operands ...int) int {
	ins := code.Make(op, operands...)
	pos := c.addInstruction(ins)

	c.setLastInstruction(op, pos)

	return pos
}

func (c *Compiler) addInstruction(ins []byte) int {
	posNewInstruction := len(c.instructions)
	c.instructions = append(c.instructions, ins...)
	return posNewInstruction
}

func (c *Compiler) setLastInstruction(op code.Opcode, pos int) {
	previous := c.lastInstruction                         // The previous instruction is set from the last instruction of the instruction passed
	last := EmittedInstruction{Opcode: op, Position: pos} // The instruction passed now becomes the last instruction of the current instruction

	c.previousInstruction = previous
	c.lastInstruction = last
}

func (c *Compiler) lastInstructionIsPop() bool {
	return c.lastInstruction.Opcode == code.OpPop
}

func (c *Compiler) removeLastPop() {
	// Remove the last OpPop by rewriting instructions[] with instructions[] till the lastInstruction position
	c.instructions = c.instructions[:c.lastInstruction.Position]

	// Replace the last instruction with previous instruction, ie, the one before that to keep proper track, since we removed the actual last instruction (OpPop)
	c.lastInstruction = c.previousInstruction
}

func (c *Compiler) replaceInstruction(pos int, newInstruction []byte) {
	for i := 0; i < len(newInstruction); i++ {
		c.instructions[pos+i] = newInstruction[i] // Add the bytes from newInstruction to c.instructions after an offset pos
	}
}

// Replace the operand of an instruction. The assumption here is that we only replace instructions of the same type with same length.
func (c *Compiler) changeOperand(opPos int, operand int) {
	op := code.Opcode(c.instructions[opPos])
	newInstruction := code.Make(op, operand)

	c.replaceInstruction(opPos, newInstruction)
}
