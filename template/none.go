package template

import (
	cond "github.com/vela-ssoc/vela-cond"
	"github.com/vela-ssoc/vela-kit/lua"
)

var NaN = &None{}

type None struct{}

func (none *None) String() string                         { return "" }
func (none *None) Type() lua.LValueType                   { return lua.LTObject }
func (none *None) AssertFloat64() (float64, bool)         { return 0, false }
func (none *None) AssertString() (string, bool)           { return "", false }
func (none *None) AssertFunction() (*lua.LFunction, bool) { return nil, false }
func (none *None) Peek() lua.LValue                       { return none }

func (none *None) Field(string) string {
	return ""
}

func (none *None) Compare(_, _ string, method cond.Method) bool {
	return false
}
