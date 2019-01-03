package commands

import (
	"fmt"
	"github.com/Fall1ngStar/sls-app-builder/function"
	"github.com/Fall1ngStar/sls-app-builder/project"
	"github.com/Fall1ngStar/sls-app-builder/utils"
	"github.com/iancoleman/strcase"
	"github.com/urfave/cli"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
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
	updateConfig(p, functionName)
	addFiles(p, functionName)
	return nil
}

func updateConfig(p *project.Project, functionName string) {
	if p.Serverless.Functions == nil {
		p.Serverless.Functions = make(map[string]function.Function)
	}

	p.Serverless.Functions[strcase.ToLowerCamel(functionName)] = function.Function{
		Handler: "src/" + functionName + ".handler",
		Layers: []string{
			"arn:aws:lambda:#{AWS::Region}:#{AWS::AccountId}:layer:proj-deps-dev",
		},
	}

	p.Serverless.UpdateConfigFile("./serverless.yml")
}

func addFiles(p *project.Project, functionName string) {
	templateFunction := function.TemplateFunction{FunctionSnake: functionName}

	funcTempalte, _ := p.Box.FindString("template_function")
	utils.WriteTemplateToFile(
		funcTempalte,
		filepath.Join("src", functionName+".py"),
		templateFunction)

	funcTestTempalte, _ := p.Box.FindString("template_test_function")
	utils.WriteTemplateToFile(
		funcTestTempalte,
		filepath.Join("src", "unit_test", "test_"+functionName+".py"),
		templateFunction)

}

func CreateEnvVariable(c *cli.Context) error {
	if c.NArg() != 1 {
		return cli.NewExitError("Usage: slapp add env <variable name>", 1)
	}

	p, err := project.LoadProject()

	if err != nil {
		return err
	}

	variableName := strcase.ToScreamingSnake(c.Args()[0])
	var variableContent string
	if c.String("from-env") != "" {
		variableContent = "${" + c.String("from-env") + "}"
	}

	for _, env := range project.ENVS {
		envFilePath := filepath.Join("environ", env+".yml")
		content, err := ioutil.ReadFile(envFilePath)
		if err != nil {
			return err
		}
		var variables map[string]interface{}

		err = yaml.Unmarshal(content, &variables)
		if err != nil {
			return err
		}

		variables[variableName] = variableContent
		result, err := yaml.Marshal(variables)
		if err != nil {
			return err
		}
		err = ioutil.WriteFile(envFilePath, result, os.FileMode(os.O_RDWR|os.O_CREATE|os.O_TRUNC))
		if err != nil {
			fmt.Println(err)
		}

	}
	provider := &p.Serverless.Provider
	if provider.Environment == nil {
		provider.Environment = make(map[string]string)
	}
	provider.Environment[variableName] = "${env:" + variableName + "}"
	p.Serverless.UpdateConfigFile("serverless.yml")
	return nil
}
