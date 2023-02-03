package template

import (
	"github.com/vela-ssoc/vela-kit/vela"
	"github.com/vela-ssoc/vela-kit/lua"
)

var xEnv vela.Environment

func NewL(L *lua.LState) *Template {
	return &Template{
		co: xEnv.Clone(L),
	}
}

func WithEnv(env vela.Environment) {
	xEnv = env
}
