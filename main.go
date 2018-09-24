package main

import (
	"fmt"
	"github.com/Fall1ngStar/sls-app-builder/project"
	"github.com/urfave/cli"
	"log"
	"os"
	"strings"
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
				fmt.Println(p.Serverless.Service)
				fmt.Println(p.Path)
				return nil
			},
		},
		{
			Name: "package",
			Action: func(c *cli.Context) error {
				proj, _ := project.LoadProject()
				head, _ := proj.Repository.Head()
				name := head.Name().Short()
				name = "5-edu"
				name = strings.Replace(strings.Title(name), "_", "-", -1)
				result := strings.Split(name, "-")[1:]
				fmt.Println(strings.Join(result, ""))
				//cmd := exec.Command("serverless", "package", "--verbose")
				//cmd.Stdout = os.Stdout
				//cmd.Run()
				return nil
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
