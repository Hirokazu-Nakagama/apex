package template

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/tj/hil/ast"
)

// TODO: debug level logs here

// Env returns the environment variable values or its default.
var Env = ast.Function{
	ArgTypes:     []ast.Type{ast.TypeString},
	VariadicType: ast.TypeString,
	ReturnType:   ast.TypeString,
	Variadic:     true,
	Callback: func(args []interface{}) (interface{}, error) {
		name := args[0].(string)

		s := os.Getenv(name)

		if s != "" {
			return s, nil
		}

		if len(args) > 1 {
			return args[1].(string), nil
		}

		return nil, fmt.Errorf("missing %q and no default is provided", name)
	},
}

// Shell returns shell command output.
var Shell = ast.Function{
	ArgTypes:     []ast.Type{ast.TypeString},
	VariadicType: ast.TypeString,
	ReturnType:   ast.TypeString,
	Variadic:     true,
	Callback: func(args []interface{}) (interface{}, error) {
		in := fmt.Sprintf(args[0].(string), args[1:]...)
		cmd := exec.Command("sh", "-c", in)

		b, err := cmd.CombinedOutput()
		if err != nil {
			return nil, fmt.Errorf("%s: %s", err, strings.Trim(string(b), " \n"))
		}

		return string(b), nil
	},
}
