package template

import (
	auxlib2 "github.com/vela-ssoc/vela-kit/auxlib"
	"github.com/vela-ssoc/vela-kit/lua"
	"github.com/vela-ssoc/vela-kit/xreflect"
)

func (ctx *Context) String() string                         { return "" }
func (ctx *Context) Type() lua.LValueType                   { return lua.LTObject }
func (ctx *Context) AssertFloat64() (float64, bool)         { return 0, false }
func (ctx *Context) AssertString() (string, bool)           { return "", false }
func (ctx *Context) AssertFunction() (*lua.LFunction, bool) { return nil, false }
func (ctx *Context) Peek() lua.LValue                       { return ctx }

func (ctx *Context) vL(L *lua.LState) int {
	val := L.Get(1)
	if val.Type() == lua.LTNil {
		return 0
	}

	ctx.keyword = append(ctx.keyword, val)
	return 0
}

func (ctx *Context) happyL(L *lua.LState) int {
	ctx.hit = true
	return 0
}

func (ctx *Context) payloadL(L *lua.LState) int {
	chunk := auxlib2.S2B(auxlib2.Format(L, 0))
	if len(ctx.payload) != 0 {
		ctx.payload = append(ctx.payload, '\n')
	}
	ctx.payload = append(ctx.payload, chunk...)
	return 0
}

func (ctx *Context) matchL(L *lua.LState) int {
	if ctx.match == nil {
		return 0
	}
	lv := L.Get(1)

	switch lv.Type() {
	case lua.LTNil:
		L.Push(lua.LFalse)
		return 1

	case lua.LTObject:
		c, ok := lv.(*Context)
		if ok {
			L.RaiseError("%p context cyclic assignment", c)
			L.Push(lua.LFalse)
			return 1
		}
	}

	L.Push(lua.LBool(ctx.match(lv)))
	return 1
}

func (ctx *Context) Index(L *lua.LState, key string) lua.LValue {
	switch key {
	case "id":
		return lua.S2L(ctx.ID)
	case "file":
		return lua.S2L(ctx.File)
	case "description":
		return lua.S2L(ctx.info.Description)
	case "raw_payload":
		return lua.B2L(ctx.Payload())
	case "raw_keyword":
		return lua.B2L(ctx.Payload())
	case "raw_tags":
		return lua.S2L(ctx.info.Tags.String())
	case "data_type":
		return lua.S2L(ctx.DataType())
	case "keyword":
		return lua.NewFunction(ctx.vL)
	case "happy":
		return lua.NewFunction(ctx.happyL)
	case "feedback":
		return lua.LBool(ctx.hit)
	case "payload":
		return lua.NewFunction(ctx.payloadL)
	case "value":
		return xreflect.ToLValue(ctx.data, L)
	case "match":
		return lua.NewFunction(ctx.matchL)
	}

	return lua.LNil
}
