package vm

import (
	"Compiler/c-monkey-v6/src/code"
	"Compiler/c-monkey-v6/src/object"
)

type Frame struct {
	fn *object.CompiledFunction
	ip int // The instruction pointer within this particular frame
}

func NewFrame(fn *object.CompiledFunction) *Frame {
	return &Frame{fn: fn, ip: -1}
}

func (f *Frame) Instructions() code.Instructions {
	return f.fn.Instructions
}
