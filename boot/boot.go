package boot

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/mitchellh/go-wordwrap"
	"github.com/tj/go-prompt"

	"github.com/apex/apex/boot/boilerplate"
)

var modulesCommand = `
  terraform get
`

var stateCommand = `
  terraform remote config \
    -backend=s3 \
    -backend-config="bucket=%s" \
    -backend-config="key=terraform/state/%s"
`

var projectConfig = `
{
  "name": "%s",
  "description": "%s",
  "memory": 128,
  "timeout": 5,
  "role": "arn:aws:iam::%s:role/lambda",
  "environment": {}
}`

// All bootstraps a project and infrastructure.
func All(region string) error {
	if err := Project(); err != nil {
		return err
	}

	help(`Would you like to manage infrastructure with Terraform?`)
	if prompt.Confirm("Use Terraform (yes/no)? ") {
		fmt.Println()
		if err := Infra(region); err != nil {
			return err
		}
	}

	help(`Setup complete :)`)

	return nil
}

// Project bootstraps a project.
func Project() error {
	help(`Enter the name of your project. It should be machine-friendly, as this is used to prefix your functions in Lambda.`)
	name := prompt.StringRequired("  Project name: ")

	help(`Enter an optional description of your project.`)
	description := prompt.String("  Project description: ")

	// TODO(tj): once we have TF -> Apex this can be removed, it's used
	// to reference the Role for now.
	help(`Enter your AWS account ID.`)
	accountID := prompt.StringRequired("  AWS account id: ")
	fmt.Println()

	logf("creating ./project.json")
	project := fmt.Sprintf(projectConfig, name, description, accountID)
	return ioutil.WriteFile("project.json", []byte(project), 0644)
}

// Infra bootstraps terraform for infrastructure management.
func Infra(region string) error {
	// TODO(tj): derp, required by TF right now?
	os.Setenv("AWS_DEFAULT_REGION", region)

	if _, err := exec.LookPath("terraform"); err != nil {
		return fmt.Errorf("terraform is not installed")
	}

	logf("creating ./infrastructure")
	if err := boilerplate.RestoreAssets(".", "infrastructure"); err != nil {
		return err
	}

	logf("creating ./functions")
	if err := boilerplate.RestoreAssets(".", "functions"); err != nil {
		return err
	}

	// TODO(tj): verify access first?
	help(`Enter the S3 bucket name for managing Terraform state.`)
	bucket := prompt.StringRequired("S3 bucket name: ")
	fmt.Println()

	if err := InfraEnv(region, bucket, "prod"); err != nil {
		return err
	}

	return InfraEnv(region, bucket, "stage")
}

// InfraEnv sets up a single infrastructure environment.
func InfraEnv(region, bucket, env string) error {
	if err := setupState(region, bucket, env); err != nil {
		return err
	}

	return setupModules(env)
}

// setupModules performs a `terraform get` for `env`.
func setupModules(env string) error {
	logf("fetching %q modules", env)
	dir := filepath.Join("infrastructure", env)
	return shell(modulesCommand, dir)
}

// setupState performs a `terraform remote config` for `env`.
func setupState(region, bucket, env string) error {
	logf("setting up %q state in bucket %q", env, bucket)
	cmd := fmt.Sprintf(stateCommand, bucket, env)
	dir := filepath.Join("infrastructure", env)
	return shell(cmd, dir)
}

// shell executes `command` in a shell within `dir`.
func shell(command, dir string) error {
	cmd := exec.Command("sh", "-c", command)
	cmd.Dir = dir

	out, err := cmd.CombinedOutput()
	if err != nil {
		// TODO(tj): make it look nice
		return fmt.Errorf("error executing command: %s: %s", out, err)
	}

	return nil
}

// help string output.
func help(s string) {
	os.Stdout.WriteString("\n")
	// TODO(tj): indent
	os.Stdout.WriteString(wordwrap.WrapString(s, 70))
	os.Stdout.WriteString("\n\n")
}

// logf outputs a log message.
func logf(s string, v ...interface{}) {
	fmt.Printf("  \033[34m[+]\033[0m %s\n", fmt.Sprintf(s, v...))
}
