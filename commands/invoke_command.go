package commands

import (
	"github.com/Fall1ngStar/sls-app-builder/project"
	"github.com/Fall1ngStar/sls-app-builder/utils"
	"github.com/iancoleman/strcase"
	"github.com/urfave/cli"
)

func InvokeFunction(c *cli.Context) error {
	if c.NArg() != 1 {
		return cli.NewExitError("Usage: "+c.App.Name+" invoke <function name>", 1)
	}

	p, err := project.LoadProject()
	if err != nil {
		return err
	}

	stage := c.String("stage")
	branch := p.GetBranchName()

	err = loadEnvFile(stage, false)
	if err != nil {
		return err
	}
	if branch == "" {
		branch = stage
	}
	err = getDepLayerVersion(p, branch)
	if err != nil {
		return err
	}

	functionName := strcase.ToLowerCamel(c.Args()[0])
	return utils.RunWithStdout("serverless", "invoke",
		"--stage", branch,
		"-f", functionName)
}
