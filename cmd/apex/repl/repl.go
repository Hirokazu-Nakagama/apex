// Package repl implements a Read-Evaluate-Print-Loop which executes commands in Lambda,
// primarily useful testing, and inspecting the Lambda environment.
package repl

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/chzyer/readline"
	"github.com/spf13/cobra"

	"github.com/apex/apex/cmd/apex/root"
)

// example output.
const example = `  Start a REPL
  $ apex repl`

// Command config.
var Command = &cobra.Command{
	Use:     "repl",
	Short:   "Interactive Lambda REPL",
	Example: example,
	RunE:    run,
}

// Initialize.
func init() {
	root.Register(Command)
}

// source for the lambda function.
var source = `
package main

import (
	"encoding/json"
	"os/exec"

	"github.com/apex/go-apex"
	"github.com/apex/log"
	"github.com/apex/log/handlers/logfmt"
)

type message struct {
	Command string
}

func init() {
	log.SetHandler(logfmt.Default)
}

func main() {
	log.Info("starting")

	apex.HandleFunc(func(event json.RawMessage, ctx *apex.Context) (interface{}, error) {
		var msg message

		if err := json.Unmarshal(event, &msg); err != nil {
			return nil, err
		}

		log.WithField("command", msg.Command).Info("exec")

		cmd := exec.Command("sh", "-c", msg.Command)
		out, err := cmd.CombinedOutput()
		return string(out), err
	})
}
`

var config = `
{
  "description": "Apex generated REPL function",
  "runtime": "golang"
}
`

type event struct {
	Command string
}

// Run command.
func run(c *cobra.Command, args []string) error {
	path := filepath.Join(os.TempDir(), "__apex_repl__")

	// TODO: flag that it has been deployed

	if err := os.MkdirAll(path, 0755); err != nil {
		return err
	}

	if err := ioutil.WriteFile(filepath.Join(path, "main.go"), []byte(source), 0755); err != nil {
		return err
	}

	if err := ioutil.WriteFile(filepath.Join(path, "function.json"), []byte(config), 0755); err != nil {
		return err
	}

	fn, err := root.Project.LoadFunctionPath("repl", path)
	if err != nil {
		return err
	}

	if err := fn.Deploy(); err != nil {
		return err
	}

	rl, err := readline.New("\033[34mlambda>\033[0m ")
	if err != nil {
		return err
	}
	defer rl.Close()

	for {
		line, err := rl.Readline()
		if err != nil {
			break
		}

		reply, _, err := fn.Invoke(event{line}, nil)
		if err != nil {
			return err
		}

		var s string

		if err := json.NewDecoder(reply).Decode(&s); err != nil {
			return err
		}

		os.Stdout.WriteString(s)
	}

	return nil
}
