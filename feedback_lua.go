package engine

import (
	"github.com/vela-ssoc/vela-kit/audit"
	"github.com/vela-ssoc/vela-kit/kind"
	"github.com/vela-ssoc/vela-kit/pipe"
	"github.com/vela-ssoc/vela-kit/lua"
)

func (fb *Feedback) String() string                         { return "" }
func (fb *Feedback) Type() lua.LValueType                   { return lua.LTObject }
func (fb *Feedback) AssertFloat64() (float64, bool)         { return 0, false }
func (fb *Feedback) AssertString() (string, bool)           { return "", false }
func (fb *Feedback) AssertFunction() (*lua.LFunction, bool) { return nil, false }
func (fb *Feedback) Peek() lua.LValue                       { return fb }

func (fb *Feedback) debugL(L *lua.LState) int {
	n := len(fb.Value)
	if n == 0 {
		L.Push(lua.S2L("[]"))
		return 1
	}

	enc := kind.NewJsonEncoder()
	enc.Arr("")
	for i := 0; i < n; i++ {
		enc.Tab("")
		ctx := fb.Value[i]
		info := ctx.Info()
		enc.KV("id", ctx.ID)
		enc.KV("file", ctx.File)
		enc.KV("description", info.Description)
		enc.KV("tags", info.Tags.String())
		enc.KV("data_type", ctx.DataType())
		enc.KV("payload", ctx.Payload())
		enc.End("},")
	}

	enc.End("]")

	L.Push(lua.B2L(enc.Bytes()))
	return 1
}

func (fb *Feedback) pipeL(L *lua.LState) int {
	n := len(fb.Value)
	if n == 0 {
		return 0
	}
	pip := pipe.NewByLua(L)
	for i := 0; i < n; i++ {
		err := pip.Call(fb.co, fb.Value[i])
		if err != nil {
			audit.Debug("engine feedback pipe call fail %v", err)
			continue
		}
	}
	return 0
}
func (fb *Feedback) Index(L *lua.LState, key string) lua.LValue {
	switch key {
	case "collect":
		return lua.GoFuncErr(fb.collect)
	case "debug":
		return lua.NewFunction(fb.debugL)
	case "pipe":
		return lua.NewFunction(fb.pipeL)
	}

	return lua.LNil
}

func NewFeedbackL(L *lua.LState) int {
	fb := NewFeedback()
	fb.co = xEnv.Clone(L)
	L.Push(fb)
	return 1
}

func CheckFeedback(L *lua.LState, idx int) *Feedback {
	val := L.Get(idx)
	if val.Type() != lua.LTObject {
		L.RaiseError("#%d invalid engine feedback object", idx)
		return nil
	}

	fb, ok := val.(*Feedback)
	if ok {
		return fb
	}
	L.RaiseError("#%d invalid engine feedback", idx)
	return nil
}
