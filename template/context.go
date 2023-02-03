package template

import (
	"fmt"
	cond "github.com/vela-ssoc/vela-cond"
	"github.com/vela-ssoc/vela-engine/header"
	"github.com/vela-ssoc/vela-kit/lua"
	"reflect"
)

type Context struct {
	ID      string
	File    string
	match   func(lv interface{}) bool
	keyword []lua.LValue
	info    *header.Info
	payload []byte
	data    interface{}
	co      *lua.LState
	hit     bool
}

func (ctx *Context) Info() *header.Info {
	return ctx.info
}

func (ctx *Context) DataType() string {
	return reflect.TypeOf(ctx.data).String()
}

func (ctx *Context) Payload() []byte {
	return ctx.payload
}

func (ctx *Context) pay(id int, v string) {
	ctx.keyword = append(ctx.keyword, lua.S2L(v))
}

func (ctx *Context) Hit() bool {
	return ctx.hit
}

func (ctx *Context) CompareKeyword(val string, method cond.Method) bool {
	n := len(ctx.keyword)
	if n == 0 {
		return false
	}

	for i := 0; i < n; i++ {
		if method(ctx.keyword[i].String(), val) {
			return true
		}
	}
	return false
}

func (ctx *Context) Compare(typ, val string, method cond.Method) bool {
	switch typ {
	case "keyword":
		return ctx.CompareKeyword(val, method)
	case "feedback":
		return lua.LBool(ctx.hit).String() == val
	}

	var peek cond.Peek

	switch item := ctx.data.(type) {
	case nil:
		return false

	case cond.Peek:
		peek = item
	case cond.FieldEx:
		peek = item.Field
	case cond.CompareEx:
		return item.Compare(typ, val, method)

	case string:
		peek = cond.String(item)

	case []byte:
		peek = cond.String(string(item))

	case func() string:
		peek = func(string) string {
			return item()
		}

	case lua.IndexEx:
		peek = func(key string) string {
			return item.Index(nil, key).String()
		}

	case lua.MetaEx:
		peek = func(key string) string {
			return item.Meta(nil, lua.S2L(key)).String()
		}

	case lua.MetaTableEx:
		peek = func(key string) string {
			return item.MetaTable(nil, key).String()
		}

	case *lua.LTable:
		peek = func(key string) string {
			return item.RawGetString(key).String()
		}

	case fmt.Stringer:
		peek = cond.String(item.String())

	default:
		return false
	}

	value := peek(typ)

	return method(value, val)
}
