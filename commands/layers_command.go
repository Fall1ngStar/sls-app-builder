package commands

import (
	"github.com/Fall1ngStar/sls-app-builder/project"
	"github.com/Fall1ngStar/sls-app-builder/utils"
	"github.com/urfave/cli"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
)

func DeployLayers(c *cli.Context) error {
	os.RemoveAll("tmp")
	err := os.MkdirAll("tmp/dist", 0777)
	if err != nil {
		return err
	}
	//defer os.RemoveAll("tmp")

	p, err := project.LoadProject()
	if err != nil {
		return err
	}
	cmd := exec.Command("pipenv", "lock", "-r")
	result, err := cmd.Output()
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

	cmd = exec.Command("pipenv", "run", "pip", "install", "-r",
		filepath.Join("tmp", "requirements.txt"), "--target",
		filepath.Join("tmp", "dist"))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		return err
	}
	err = os.Chdir("tmp")
	if err != nil {
		return err
	}

	cmd = exec.Command("serverless", "deploy")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		return err
	}
	return nil
}
