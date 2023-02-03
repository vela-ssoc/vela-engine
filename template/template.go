package template

import (
	"bytes"
	"fmt"
	cond "github.com/vela-ssoc/vela-cond"
	"github.com/vela-ssoc/vela-engine/header"
	"github.com/vela-ssoc/vela-engine/match"
	"github.com/vela-ssoc/vela-kit/execpt"
	"github.com/vela-ssoc/vela-kit/lua"
	"github.com/vela-ssoc/vela-kit/xreflect"
	"gopkg.in/yaml.v3"
	"io"
	"os"
)

type Template struct {
	ID           string           `yaml:"id"`
	Info         *header.Info     `yaml:"info"`
	Condition    *match.Condition `yaml:"condition"`
	BeforeScript string           `yaml:"before"`
	AfterScript  string           `yaml:"after"`

	Filename string         `yaml:"-"`
	before   *lua.LFunction `yaml:"-"`
	after    *lua.LFunction `yaml:"-"`
	co       *lua.LState    `yaml:"-"`
}

func (tmpl *Template) context(v interface{}, co *lua.LState) *Context {
	ctx := &Context{
		ID:   tmpl.ID,
		co:   co,
		info: tmpl.Info,
		data: v,
		File: tmpl.Filename,
	}

	ctx.match = func(lv interface{}) bool {
		return tmpl.Condition.Match(lv,
			cond.Payload(ctx.pay),
			cond.WithCo(tmpl.co))
	}

	return ctx
}

func (tmpl *Template) PrepareBefore() error {
	if len(tmpl.BeforeScript) == 0 {
		return nil
	}

	before, err := tmpl.Prepare(tmpl.BeforeScript)
	if err != nil {
		return err
	}

	tmpl.before = before
	return nil
}

func (tmpl *Template) PrepareAfter() error {
	if len(tmpl.AfterScript) == 0 {
		return nil
	}

	after, err := tmpl.Prepare(tmpl.AfterScript)
	if err != nil {
		return err
	}
	tmpl.after = after
	return nil
}

func (tmpl *Template) ReadFile(path string) error {
	tmpl.Filename = path

	r, err := os.Open(path)
	if err != nil {
		return err
	}
	defer r.Close()

	return tmpl.Decoder(r)
}

func (tmpl *Template) Decoder(r io.Reader) error {
	decode := yaml.NewDecoder(r)
	err := decode.Decode(tmpl)
	if err != nil {
		return err
	}

	err = tmpl.PrepareBefore()
	if err != nil {
		return err
	}

	err = tmpl.PrepareAfter()
	if err != nil {
		return err
	}

	err = tmpl.Condition.Prepare()

	return err
}

func (tmpl *Template) Call(v interface{}) *Context {
	errs := execpt.New()
	ctx := tmpl.context(v, tmpl.co)
	val := xreflect.ToLValue(v, tmpl.co)

	if tmpl.before != nil {
		errs.Try("before", xEnv.Call(tmpl.Coroutine(), tmpl.before, ctx, val))
	}

	if tmpl.Condition.Match(ctx, cond.Payload(ctx.pay), cond.WithCo(tmpl.co)) {
		ctx.hit = true
	}

	if tmpl.after != nil {
		errs.Try("after", xEnv.Call(tmpl.Coroutine(), tmpl.after, ctx, val))
	}

	return ctx
}

func (tmpl *Template) Coroutine() *lua.LState {
	if tmpl.co == nil {
		tmpl.co = xEnv.Coroutine()
	}
	return tmpl.co
}

func (tmpl *Template) Prepare(script string) (*lua.LFunction, error) {
	var chunk bytes.Buffer
	chunk.WriteString("return function(ctx)\n")
	chunk.WriteString(script)
	chunk.WriteString("\nend")

	co := tmpl.Coroutine()
	err := xEnv.DoString(co, chunk.String())
	if err != nil {
		return nil, err
	}

	top := co.Get(-1)
	if top.Type() != lua.LTFunction {
		return nil, fmt.Errorf("script compile fail %v", err)
	}

	fn := top.(*lua.LFunction)
	co.SetTop(0)
	return fn, nil
}
