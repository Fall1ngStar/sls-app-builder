package utils

import (
	"github.com/urfave/cli"
	"os/exec"
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
	requiredExecutables := []string {
		"npm",
		"serverless",
		"pipenv",
		"aws",
	}
	for _, executable := range requiredExecutables {
		if !checkExecutableInPath(executable) {
			return cli.NewExitError("Could not find \"" + executable + "\" in path", 1)
		}
	}
	return nil
}