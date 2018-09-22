package main

import (
	"github.com/fall1ngstar/builder/project"
	"github.com/urfave/cli"
	"gopkg.in/libgit2/git2go.v26"
	"log"
	"os"
	"fmt"
)

func createProject(path string) {
	os.Mkdir(path, 0644)
	git.InitRepository(path, false)
}

func main() {
	app := cli.NewApp()
	app.Name = "Builder"
	app.Usage = "Building serverless app"

	app.Commands = []cli.Command{
		{
			Name:   "create",
			Usage:  "Create a project folder",
			Action: project.CreateProject,
		},
		{
			Name: "check",
			Usage: "Check requirements for project creation",
			Action: func (c *cli.Context) error {
				fmt.Println("Serverless:", project.CheckExecutableInPath("serverless"))
				fmt.Println("Python:", project.CheckExecutableInPath("python"))
				fmt.Println("NPM:", project.CheckExecutableInPath("npm"))
				return nil
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
