package code

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type Instructions []byte // Instructions can be any length of bytes
type Opcode byte         // Opcodes are always a single byte

const ( // Each Opcode is going to have a readable name instead of some arbitrary byte value that will need to be memorized
	OpConstant Opcode = iota // OpConstant is for the operation constants , ie, operands
	OpAdd
	OpSub
	OpMul
	OpDiv
	OpPop
	OpTrue
	OpFalse
	OpEqual
	OpNotEqual
	OpGreaterThan
	OpGreaterThanOrEqual
	OpMinus
	OpBang
	OpJumpNotTruthy // Value on top of the stack is not truthy, aka, it is false
	OpJump          // Unconditional jump
	OpNull
	OpGetGlobal
	OpSetGlobal
)

type Definition struct { // To keep track of how many operands an opcode has and make it more readable
	Name          string
	OperandWidths []int
}

var definitions = map[Opcode]*Definition{
	// OpConstant operand is 2 bytes wide, ie, uint16. An array is used to indicate that the length of the array (1):
	// there is only one operand and the value of that element (2) is the width
	OpConstant:           {"OpConstant", []int{2}},
	OpAdd:                {"OpAdd", []int{}},
	OpSub:                {"OpSub", []int{}},
	OpMul:                {"OpMul", []int{}},
	OpDiv:                {"OpDiv", []int{}},
	OpPop:                {"OpPop", []int{}},
	OpTrue:               {"OpTrue", []int{}},
	OpFalse:              {"OpFalse", []int{}},
	OpEqual:              {"OpEqual", []int{}},
	OpNotEqual:           {"OpNotEqual", []int{}},
	OpGreaterThan:        {"OpGreaterThan", []int{}},
	OpGreaterThanOrEqual: {"OpGreaterthanOrEqual", []int{}},
	OpMinus:              {"OpMinus", []int{}},
	OpBang:               {"OpBang", []int{}},
	OpJumpNotTruthy:      {"OpJumpNotTruthy", []int{2}},
	OpJump:               {"OpJump", []int{2}},
	OpNull:               {"OpNull", []int{}},
	OpGetGlobal:          {"OpGetGlobal", []int{2}},
	OpSetGlobal:          {"OpSetGlobal", []int{2}},
}

func Lookup(op byte) (*Definition, error) {
	def, ok := definitions[Opcode(op)]
	if !ok {
		return nil, fmt.Errorf("opcode %d undefined", op)
	}

	return def, nil
}

func Make(op Opcode, operands ...int) []byte {
	def, ok := definitions[op]
	if !ok {
		return []byte{}
	}

	instructionLen := 1
	for _, w := range def.OperandWidths {
		instructionLen += w
	}

	instruction := make([]byte, instructionLen)
	instruction[0] = byte(op)

	offset := 1
	for i, o := range operands {
		width := def.OperandWidths[i]
		switch width {
		case 2:
			binary.BigEndian.PutUint16(instruction[offset:], uint16(o))
		}
		offset += width
	}

	return instruction
}

func (ins Instructions) String() string {
	var out bytes.Buffer

	i := 0
	for i < len(ins) {
		def, err := Lookup(ins[i])
		if err != nil {
			fmt.Fprintf(&out, "ERROR: %s\n", err)
			continue
		}

		operands, read := ReadOperands(def, ins[i+1:])
		fmt.Fprintf(&out, "%04d %s\n", i, ins.fmtInstruction(def, operands))
		i += 1 + read
	}

	return out.String()
}

func (ins Instructions) fmtInstruction(def *Definition, operands []int) string {
	operandCount := len(def.OperandWidths)

	if len(operands) != operandCount {
		return fmt.Sprintf("ERROR: operand len %d does not match defined %d\n", len(operands), operandCount)
	}

	switch operandCount {
	case 0:
		return def.Name
	case 1:
		return fmt.Sprintf("%s %d", def.Name, operands[0])
	}

	return fmt.Sprintf("ERROR: unhandled operandCount for %s\n", def.Name)
}

func ReadOperands(def *Definition, ins Instructions) ([]int, int) {
	operands := make([]int, len(def.OperandWidths))
	offset := 0

	for i, width := range def.OperandWidths {
		switch width {
		case 2:
			operands[i] = int(ReadUint16(ins[offset:]))
		}

		offset += width
	}

	return operands, offset
}

func ReadUint16(ins Instructions) uint16 {
	return binary.BigEndian.Uint16(ins)
}
