package main

import (
	"github.com/fall1ngstar/builder/project"
	"github.com/urfave/cli"
	"gopkg.in/libgit2/git2go.v26"
	"log"
	"os"
)

func createProject(path string) {
	os.Mkdir(path, 0644)
	git.InitRepository(path, false)
}

func main() {
	project.CheckServerlessRequirements()
	app := cli.NewApp()
	app.Name = "Builder"
	app.Usage = "Building serverless app"

	app.Commands = []cli.Command{
		{
			Name:   "create",
			Usage:  "Create a project folder",
			Action: project.CreateProject,
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
