package main

import (
	"github.com/Fall1ngStar/sls-app-builder/commands"
	"github.com/Fall1ngStar/sls-app-builder/project"
	"github.com/Fall1ngStar/sls-app-builder/utils"
	"github.com/urfave/cli"
	"log"
	"os"
)

func main() {
	app := cli.NewApp()
	app.Name = "SLApp"
	app.Usage = "A CLI to create ServerLess Applications in Python"
	app.Version = "0.0.1"
	app.Commands = []cli.Command{
		{
			Name:   "create",
			Usage:  "Create a project folder",
			Action: project.CreateProject,
		},
		{
			Name:   "check",
			Usage:  "Check executables requirements for project creation",
			Action: utils.CheckRequiredExecutables,
		},
		{
			Name:   "layers",
			Usage:  "Deploy the dependency layer for this environment",
			Action: commands.DeployLayers,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "stage, s",
					Value: "DEV",
					Usage: "Stage to deploy",
				},
				cli.StringFlag{
					Name: "env",
				},
			},
		},
		{
			Name:  "deploy",
			Usage: "Deploy the current projet",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "stage, s",
					Value: "DEV",
					Usage: "Stage to deploy",
				},
				cli.BoolFlag{
					Name: "gitlab, g",
				},
				cli.StringFlag{
					Name:  "function, f",
					Usage: "Deploy a specific function",
				},
			},
			Action: commands.DeployProject,
		},
		{
			Name: "add",
			Subcommands: []cli.Command{
				{
					Name:   "function",
					Action: commands.CreateFunction,
				},
				{
					Name:   "env",
					Action: commands.CreateEnvVariable,
					Flags: []cli.Flag{
						cli.StringFlag{
							Name: "from-env, e",
						},
					},
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
