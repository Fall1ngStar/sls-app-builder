package utils

import (
	"github.com/urfave/cli"
	"os"
	"os/exec"
	"text/template"
)

func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func checkExecutableInPath(executable string) bool {
	_, err := exec.LookPath(executable)
	return err == nil
}

func CheckRequiredExecutables(c *cli.Context) error {
	requiredExecutables := []string{
		"npm",
		"serverless",
		"pipenv",
		"aws",
	}
	for _, executable := range requiredExecutables {
		if !checkExecutableInPath(executable) {
			return cli.NewExitError("Could not find \""+executable+"\" in path", 1)
		}
	}
	return nil
}

func WriteTemplateToFile(tmpl, path string, content interface{}) {
	t, _ := template.New("tmp").Parse(tmpl)
	file, _ := os.Create(path)
	defer file.Close()
	t.Execute(file, content)
}

func RunWithStdout(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func RunWithOutput(name string, args ...string) ([]byte, error) {
	cmd := exec.Command(name, args...)
	return cmd.Output()
}