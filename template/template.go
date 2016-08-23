// Package template provides Terraform-like template interpolation and functions.
package template

import (
	"github.com/tj/hil"
	"github.com/tj/hil/ast"
)

// Template wraps HIL to provide a slightly higher level interface.
type Template struct {
	vars  map[string]ast.Variable
	funcs map[string]ast.Function
}

// New template.
func New() *Template {
	return &Template{
		vars:  make(map[string]ast.Variable),
		funcs: make(map[string]ast.Function),
	}
}

// AddFunction adds function `name` to the global scope.
func (t *Template) AddFunction(name string, value ast.Function) {
	t.funcs[name] = value
}

// AddVariable adds variable `name` to the global scope.
func (t *Template) AddVariable(name string, value interface{}) {
	val, err := hil.InterfaceToVariable(value)
	if err != nil {
		panic(err)
	}

	t.vars[name] = val
}

// Eval returns the evaluated source.
func (t *Template) Eval(s string) (string, error) {
	node, err := hil.Parse(s)
	if err != nil {
		return "", err
	}

	config := &hil.EvalConfig{
		GlobalScope: &ast.BasicScope{
			VarMap:  t.vars,
			FuncMap: t.funcs,
		},
	}

	out, err := hil.Eval(node, config)
	if err != nil {
		return "", err
	}

	return out.Value.(string), nil
}
