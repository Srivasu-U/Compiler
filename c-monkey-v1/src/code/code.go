package code

import (
	"encoding/binary"
	"fmt"
)

type Instructions []byte // Instructions can be any length of bytes
type Opcode byte         // Opcodes are always a single byte

const ( // Each Opcode is going to have a readable name instead of some arbitrary byte value that will need to be memorized
	OpConstant Opcode = iota // OpConstant is for the operation constants , ie, operands
)

type Definition struct { // To keep track of how many operands an opcode has and make it more readable
	Name          string
	OperandWidths []int
}

var definitions = map[Opcode]*Definition{
	OpConstant: {"OpConstant", []int{2}}, // OpConstant operand is 2 bytes wide, ie, uint16
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
