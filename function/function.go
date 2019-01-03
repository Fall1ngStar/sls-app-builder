package function

import (
    "github.com/Fall1ngStar/sls-app-builder/project"
    "github.com/Fall1ngStar/sls-app-builder/utils"
    "github.com/iancoleman/strcase"
    "path/filepath"
)

type Function struct {
    Handler string
    //Events  []Event `yaml:",omitempty"`
}

type Event struct {

}


type TemplateFunction struct {
     FunctionSnake string
}

func UpdateConfig(p *project.Project, functionName string) {
    if p.Serverless.Functions == nil {
        p.Serverless.Functions = make(map[string]Function)
    }

    p.Serverless.Functions[strcase.ToLowerCamel(functionName)] = Function{
        Handler: "src/" + functionName + ".handler",
    }

    p.Serverless.UpdateConfigFile("./serverless.yml")
}

func AddFiles(p *project.Project, functionName string) {
    templateFunction := TemplateFunction{FunctionSnake: functionName}

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
