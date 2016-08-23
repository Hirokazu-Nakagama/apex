package template

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEnv_present(t *testing.T) {
	tmpl := New()
	tmpl.AddFunction("env", Env)

	os.Setenv("LOGGLY_TOKEN", "123")

	s, err := tmpl.Eval(`{ "foo": "${env('LOGGLY_TOKEN')}" }`)
	assert.NoError(t, err)
	assert.Equal(t, `{ "foo": "123" }`, s)
}

func TestEnv_fallback(t *testing.T) {
	tmpl := New()
	tmpl.AddFunction("env", Env)

	os.Setenv("LOGGLY_TOKEN", "")

	s, err := tmpl.Eval(`{ "foo": "${env('LOGGLY_TOKEN', 'oh no')}" }`)
	assert.NoError(t, err)
	assert.Equal(t, `{ "foo": "oh no" }`, s)
}

func TestEnv_missing(t *testing.T) {
	tmpl := New()
	tmpl.AddFunction("env", Env)

	os.Setenv("LOGGLY_TOKEN", "")

	_, err := tmpl.Eval(`{ "foo": "${env('LOGGLY_TOKEN')}" }`)
	assert.EqualError(t, err, `env: missing "LOGGLY_TOKEN" and no default is provided`)
}

func TestShell_ok(t *testing.T) {
	tmpl := New()
	tmpl.AddFunction("shell", Shell)

	s, err := tmpl.Eval(`{ "foo": "${shell('echo hello')}" }`)
	assert.NoError(t, err)
	assert.Equal(t, "{ \"foo\": \"hello\n\" }", s)
}

func TestShell_args(t *testing.T) {
	tmpl := New()
	tmpl.AddFunction("shell", Shell)

	s, err := tmpl.Eval(`{ "foo": "${shell('echo %s', 'hello')}" }`)
	assert.NoError(t, err)
	assert.Equal(t, "{ \"foo\": \"hello\n\" }", s)
}

func TestShell_error(t *testing.T) {
	tmpl := New()
	tmpl.AddFunction("shell", Shell)

	_, err := tmpl.Eval(`{ "foo": "${shell('asdf')}" }`)
	assert.EqualError(t, err, `shell: exit status 127: sh: asdf: command not found`)
}
