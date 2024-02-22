package engine

import (
	"github.com/vela-ssoc/vela-engine/template"
	"github.com/vela-ssoc/vela-kit/exception"
	"github.com/vela-ssoc/vela-kit/lua"
	"github.com/vela-ssoc/vela-kit/vela"
	vswitch "github.com/vela-ssoc/vela-switch"
	"path/filepath"
)

var xEnv vela.Environment

/*
local e = vela.engine{
	rule = "rule.d/a.yml",
	tags = {"process" , "github"},
}


local fb = vela.engine.feedback()

local vela.require("process.lua")

local e = vela.engine.load("3rd/process.zip?tags=123&tags=123")
vela.ps().pipe(s.scan)

e.case("feedback = true").pipe(fb.collect)

vela.ps().pipe(e.scan)

engine.with(fb)

*/

func NewEngineAttachL(L *lua.LState) int {
	name := L.CheckString(1)
	info, err := xEnv.Third(name)
	if err != nil {
		L.RaiseError("%s third load fail %v", name, err)
		return 0
	}

	e := &Engine{
		co:    xEnv.Clone(L),
		vsh:   vswitch.NewL(L),
		catch: exception.New(),
	}

	if info.IsZip() {
		e.rules = []string{filepath.Join(info.File(), "*.yaml")}
	} else {
		e.rules = []string{info.File()}
	}

	L.Push(e)
	return 1
}

func NewEngineL(L *lua.LState) int {
	e := NewEngine(L)
	L.Push(e)
	return 1
}

func WithEnv(env vela.Environment) {
	xEnv = env
	template.WithEnv(env)
	kv := lua.NewUserKV()
	kv.Set("feedback", lua.NewFunction(NewFeedbackL))
	kv.Set("attach", lua.NewFunction(NewEngineAttachL))
	xEnv.Set("engine", lua.NewExport("vela.engine.export",
		lua.WithTable(kv),
		lua.WithFunc(NewEngineL)),
	)
}
