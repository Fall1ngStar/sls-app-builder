package commands

import (
    "github.com/Fall1ngStar/sls-app-builder/function"
    "github.com/Fall1ngStar/sls-app-builder/project"
    "github.com/iancoleman/strcase"
    "github.com/urfave/cli"
)

func CreateFunction(c *cli.Context) error {
    p, err := project.LoadProject()
    if err != nil {
        return err
    }

    if c.NArg() < 1 {
        return cli.NewExitError("Usage: slapp add function <function name>", 1)
    }

    functionName := c.Args()[0]

    if p.Serverless.Functions == nil {
        p.Serverless.Functions = make(map[string]function.Function)
    }

    p.Serverless.Functions[strcase.ToLowerCamel(functionName)] = function.Function{
        Handler: "src/" + functionName + ".handler",
    }

    p.Serverless.UpdateConfigFile("./serverless.yml")
    return nil
}
