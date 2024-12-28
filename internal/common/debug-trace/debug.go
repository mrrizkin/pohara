package debugtrace

import (
	"bytes"
	"runtime/debug"

	"github.com/DataDog/gostackparse"
)

type StackFrame struct {
	Function string `json:"function"`
	File     string `json:"file"`
	Line     int    `json:"line"`
}

func StackTrace() ([]StackFrame, error) {
	stack := debug.Stack()
	goroutines, err := gostackparse.Parse(bytes.NewReader(stack))
	if err != nil || len(err) > 0 {
		return nil, err[0]
	}

	frames := make([]StackFrame, 0)

	for _, goroutine := range goroutines {
		for _, stack := range goroutine.Stack {
			frames = append(frames, StackFrame{
				Function: stack.Func,
				Line:     stack.Line,
				File:     stack.File,
			})
		}
	}

	return frames, nil
}
