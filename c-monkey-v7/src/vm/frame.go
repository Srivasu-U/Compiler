package vm

import (
	"Compiler/c-monkey-v7/src/code"
	"Compiler/c-monkey-v7/src/object"
)

type Frame struct {
	fn          *object.CompiledFunction
	ip          int // The instruction pointer within this particular frame
	basePointer int // The position of the stack pointer before the execution of the function
}

func NewFrame(fn *object.CompiledFunction, basePointer int) *Frame {
	return &Frame{fn: fn, ip: -1, basePointer: basePointer}
}

func (f *Frame) Instructions() code.Instructions {
	return f.fn.Instructions
}
