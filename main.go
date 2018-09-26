package main

import (
	"fmt"
	"github.com/Fall1ngStar/sls-app-builder/project"
	"github.com/urfave/cli"
	"log"
	"os"
	"path/filepath"
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
		},
		{
			Name:  "check",
			Usage: "Check requirements for project creation",
			Action: func(c *cli.Context) error {
				fmt.Println("Serverless:", project.CheckExecutableInPath("serverless"))
				fmt.Println("Python:", project.CheckExecutableInPath("python"))
				fmt.Println("NPM:", project.CheckExecutableInPath("npm"))
				return nil
			},
		},
		{
			Name: "cfg",
			Action: func(c *cli.Context) error {
				p, err := project.LoadProject()
				if err != nil {
					return err
				}
				//key := p.Serverless["service"]
				fmt.Println(p.Serverless)
				//fmt.Println(p.Path)
				return nil
			},
		},
		{
			Name: "package",
			Action: func(c *cli.Context) error {
				proj, _ := project.LoadProject()
				fmt.Println(proj.GetBranchName())
				//cmd := exec.Command("serverless", "package", "--verbose")
				//cmd.Stdout = os.Stdout
				//cmd.Run()
				return nil
			},
		},
		{
			Name: "test",
			Action: func(c *cli.Context) error {
				p, _ := project.LoadProject()
				fmt.Println(filepath.Dir(p.Path))
				fmt.Println(filepath.Clean(p.Path))
				wd, _ := os.Getwd()
				fmt.Println(filepath.Dir(wd))
				return nil
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
