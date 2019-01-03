package main

import (
	"fmt"
	"github.com/Fall1ngStar/sls-app-builder/commands"
	"github.com/Fall1ngStar/sls-app-builder/project"
	"github.com/Fall1ngStar/sls-app-builder/utils"
	"github.com/urfave/cli"
	"log"
	"os"
)

func main() {
	app := cli.NewApp()
	app.Name = "Builder"
	app.Usage = "Building serverless app"
	app.Version = "0.0.1"
	app.Commands = []cli.Command{
		{
			Name:   "create",
			Usage:  "Create a project folder",
			Action: project.CreateProject,
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name: "skip-pipenv",
				},
			},
		},
		{
			Name:   "check",
			Usage:  "Check executables requirements for project creation",
			Action: utils.CheckRequiredExecutables,
		},
		{
			Name: "layers",
			Action: commands.DeployLayers,
		},
		{
			Name: "test",
			Action: func(c *cli.Context) error {
				os.Chdir("tmp")
				return nil
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
			},
			Action: func(c *cli.Context) error {
				fmt.Println(c.FlagNames())
				fmt.Println(c.String("stage"))
				return nil
			},
		},
		{
			Name: "add",
			Subcommands: []cli.Command{
				{
					Name:   "function",
					Action: commands.CreateFunction,
				},
				{
					Name:"env",
					Action: commands.CreateEnvVariable,
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:"from-env, e",
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
