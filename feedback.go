package engine

import (
	"github.com/vela-ssoc/vela-engine/template"
	"github.com/vela-ssoc/vela-kit/lua"
)

type Feedback struct {
	co    *lua.LState
	Value []*template.Context
}

func (fb *Feedback) collect(v ...interface{}) error {
	n := len(v)
	if n == 0 {
		return nil
	}

	ctx, ok := v[0].(*template.Context)
	if ok {
		fb.Value = append(fb.Value, ctx)
	}

	return nil
}

func (fb *Feedback) append(ctx *template.Context) {
	fb.Value = append(fb.Value, ctx)
}

func NewFeedback() *Feedback {
	return &Feedback{}
}
