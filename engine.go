package engine

import (
	"github.com/vela-ssoc/vela-engine/header"
	"github.com/vela-ssoc/vela-engine/template"
	"github.com/vela-ssoc/vela-kit/exception"
	"github.com/vela-ssoc/vela-kit/lua"
	vswitch "github.com/vela-ssoc/vela-switch"
	"path/filepath"
	"strings"
)

type Engine struct {
	catch *exception.Cause
	co    *lua.LState
	vsh   *vswitch.Switch
	tags  []string
	rules []string

	templates []*template.Template
}

func (e *Engine) MatchInfo(info *header.Info) bool {
	n := len(e.tags)
	if n == 0 {
		return true
	}

	return info.Tags.Have(e.tags...)
}

func (e *Engine) ReadRuleFile(filename string) {
	tmpl := template.NewL(e.co)
	err := tmpl.ReadFile(filename)
	if err != nil {
		e.catch.Try(filename, err)
		return
	}

	if e.MatchInfo(tmpl.Info) {
		e.templates = append(e.templates, tmpl)
	}
}

func (e *Engine) Stop(v ...interface{}) error {
	return nil
}

func (e *Engine) scan(args ...interface{}) error {
	if err := e.catch.Wrap(); err != nil {
		return err
	}

	if len(args) == 0 {
		return nil
	}

	e.Match(args[0])
	return nil
}

func (e *Engine) compile() {
	n := len(e.rules)
	if n == 0 {
		return
	}

	for i := 0; i < n; i++ {
		rule := e.rules[i]
		if strings.Index(rule, "*") == -1 {
			e.ReadRuleFile(rule)
			continue
		}

		files, err := filepath.Glob(rule)
		if err != nil {
			e.catch.Try(rule, err)
			continue
		}

		for _, file := range files {
			e.ReadRuleFile(file)
		}
	}
}

func (e *Engine) Match(v interface{}) *Feedback {
	fb := NewFeedback()

	n := len(e.templates)
	if n == 0 {
		return fb
	}

	for i := 0; i < n; i++ {
		ctx := e.templates[i].Call(v)
		e.vsh.Do(ctx)
		if ctx.Hit() {
			fb.append(ctx)
		}
	}

	return fb
}
