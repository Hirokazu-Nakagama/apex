package template

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tj/hil/ast"
)

var lowerCase = ast.Function{
	ArgTypes:   []ast.Type{ast.TypeString},
	ReturnType: ast.TypeString,
	Variadic:   false,
	Callback: func(inputs []interface{}) (interface{}, error) {
		input := inputs[0].(string)
		return strings.ToLower(input), nil
	},
}

func TestTemplate_AddVariable(t *testing.T) {
	tmpl := New()
	tmpl.AddVariable("bar", "baz")

	s, err := tmpl.Eval(`{ "foo": "${bar}" }`)
	assert.NoError(t, err)
	assert.Equal(t, `{ "foo": "baz" }`, s)
}

func TestTemplate_AddFunction(t *testing.T) {
	tmpl := New()
	tmpl.AddVariable("bar", "BAZ")
	tmpl.AddFunction("lower", lowerCase)

	s, err := tmpl.Eval(`{ "foo": "${lower(bar)}" }`)
	assert.NoError(t, err)
	assert.Equal(t, `{ "foo": "baz" }`, s)
}

func TestTemplate_variableMissing(t *testing.T) {
	tmpl := New()

	_, err := tmpl.Eval(`{
    "foo": "${bar}"
  }`)

	assert.Error(t, err, `2:16: unknown variable accessed: bar`)
}
