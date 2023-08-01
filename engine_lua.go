package engine

import (
	"fmt"
	"github.com/vela-ssoc/vela-engine/template"
	"github.com/vela-ssoc/vela-kit/auxlib"
	"github.com/vela-ssoc/vela-kit/exception"
	"github.com/vela-ssoc/vela-kit/lua"
	vswitch "github.com/vela-ssoc/vela-switch"
)

func (e *Engine) String() string                         { return "" }
func (e *Engine) Type() lua.LValueType                   { return lua.LTObject }
func (e *Engine) AssertFloat64() (float64, bool)         { return 0, false }
func (e *Engine) AssertString() (string, bool)           { return "", false }
func (e *Engine) AssertFunction() (*lua.LFunction, bool) { return nil, false }
func (e *Engine) Peek() lua.LValue                       { return e }

func (e *Engine) withL(L *lua.LState) int {
	fb := CheckFeedback(L, 1)
	n := len(e.templates)
	if n == 0 {
		return 0
	}

	for i := 0; i < n; i++ {
		ctx := e.templates[i].Call(template.NaN)
		e.vsh.Do(ctx)
		if ctx.Hit() {
			fb.append(ctx)
		}
	}
	return 0
}

func (e *Engine) tagsL(L *lua.LState) int {
	return 0
}

func (e *Engine) Index(L *lua.LState, key string) lua.LValue {
	switch key {
	case "case":
		return e.vsh.Index(L, key)
	case "stop":
		return lua.GoFuncErr(e.Stop)
	case "scan":
		return lua.GoFuncErr(e.scan)
	case "with":
		return lua.NewFunction(e.withL)
	case "tags":
		return lua.NewFunction(e.tagsL)
	case "err":
		if err := e.catch.Wrap(); err != nil {
			return lua.S2L(err.Error())
		}
	}

	return lua.LNil
}
func (e *Engine) NewIndex(L *lua.LState, key string, val lua.LValue) {
	switch key {
	case "rules":
		switch val.Type() {
		case lua.LTString:
			e.SetRules([]string{val.String()})
		case lua.LTTable:
			e.SetRules(auxlib.LTab2SS(val.(*lua.LTable)))
		}

	case "tags":
		switch val.Type() {
		case lua.LTString:
			e.tags = []string{val.String()}
		case lua.LTTable:
			e.tags = auxlib.LTab2SS(val.(*lua.LTable))
		}
	}
}

func (e *Engine) SetRules(v []string) {
	e.rules = v
}

func NewEngine(L *lua.LState) *Engine {
	e := &Engine{
		co:    xEnv.Clone(L),
		vsh:   vswitch.NewL(L),
		catch: exception.New(),
	}

	val := L.Get(1)
	switch val.Type() {
	case lua.LTString:
		e.SetRules([]string{val.String()})
		e.compile()

	case lua.LTTable:
		tab := val.(*lua.LTable)
		tab.Range(func(key string, lv lua.LValue) {
			e.NewIndex(L, key, lv)
		})

		e.compile()
	default:
		e.catch.Try("engine", fmt.Errorf("init engine invalid configure must be engine{cfg} or engine(string)"))
	}

	return e
}
