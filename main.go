package main

import (
    "fmt"
    "github.com/Fall1ngStar/sls-app-builder/commands"
    "github.com/Fall1ngStar/sls-app-builder/project"
    "github.com/Fall1ngStar/sls-app-builder/utils"
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
            Name:   "check",
            Usage:  "Check executables requirements for project creation",
            Action: utils.CheckRequiredExecutables,
        },
        {
            Name: "cfg",
            Action: func(c *cli.Context) error {
                p, err := project.LoadProject()
                if err != nil {
                    return err
                }
                fmt.Println(p.Serverless)
                return nil
            },
        },
        {
            Name: "package",
            Action: func(c *cli.Context) error {
                proj, _ := project.LoadProject()
                fmt.Println(proj.GetBranchName())
                return nil
            },
        },
        {
            Name: "test",
            Action: func(c *cli.Context) error {
                wd, _ := os.Getwd()
                fmt.Println(filepath.Dir(wd))
                fmt.Println(filepath.Ext(wd))
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
            Name:"add",
            Subcommands: []cli.Command{
                {
                  Name:"function",
                  Action: commands.CreateFunction,
                },
            },
        },
    }

    err := app.Run(os.Args)
    if err != nil {
        log.Fatal(err)
    }
}
