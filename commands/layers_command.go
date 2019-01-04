package commands

import (
	"fmt"
	"github.com/Fall1ngStar/sls-app-builder/project"
	"github.com/Fall1ngStar/sls-app-builder/utils"
	"github.com/urfave/cli"
	"io/ioutil"
	"os"
	"path/filepath"
)

func DeployLayers(c *cli.Context) error {
	// TODO check bug directory first letter removed
	err := os.MkdirAll(filepath.Join("tmp", "dist", "ppython"), 0777)
	if err != nil {
		return err
	}
	defer func() {
		err := os.RemoveAll("tmp")
		fmt.Println(err)
	}()

	p, err := project.LoadProject()
	if err != nil {
		return err
	}
	stage := c.String("stage")
	branch := p.GetBranchName()
	if branch == "" {
		branch = stage
	}

	result, err := utils.RunWithOutput("pipenv", "lock", "-r")
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(filepath.Join("tmp", "requirements.txt"), result, 0777)
	if err != nil {
		return err
	}
	layersTemplate, err := p.Box.FindString("template_serverless_layers")
	if err != nil {
		return err
	}
	utils.WriteTemplateToFile(
		layersTemplate,
		filepath.Join("tmp", "serverless.yml"),
		struct {
			ProjectName string
		}{ProjectName: p.Serverless.Service.Name})

	err = utils.RunWithStdout("pipenv", "run", "pip", "install", "-r",
		filepath.Join("tmp", "requirements.txt"), "--target",
		filepath.Join("tmp", "dist", "ppython"))
	if err != nil {
		return err
	}

	err = os.Chdir("tmp")
	if err != nil {
		return err
	}
	defer os.Chdir("..")

	err = utils.RunWithStdout("serverless", "deploy", "--stage", branch)
	if err != nil {
		return err
	}
	return nil
}
