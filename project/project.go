package project

import (
	"fmt"
	"github.com/gobuffalo/packr"
	"github.com/urfave/cli"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/yaml.v2"
	"os"
	"os/exec"
	"path"
	"text/template"
)

var ENVS = []string{
	"DEV",
	"QUAL",
	"PROD",
}

type EnvTemplate struct {
	EnvName string
}

type ServerlessConfig struct {
	Service  string
	Provider ProviderServerlessConfig `yaml:",omitempty"`
	Custom   CustomServerlessConfig   `yaml:",omitempty"`
}

type ProviderServerlessConfig struct {
	Name       string
	Region     string
	Runtime    string
	MemorySize int
}

type CustomServerlessConfig struct {
	Pip string
}

type Project struct {
	Path       string
	Repository *git.Repository
	Box        packr.Box
	Serverless ServerlessConfig
}

func CreateProject(c *cli.Context) error {
	var project Project
	if c.NArg() == 0 {
		project.Path = "."
	} else if c.NArg() == 1 {
		project.Path = c.Args()[0]
		err := project.createRootProjectFolder()
		if err != nil {
			return err
		}
	} else {
		return cli.NewExitError("Wrong argument number provided", 1)
	}
	err := project.initProjectGitRepository()
	if err != nil {
		return err
	}
	err = project.addSubFolders()
	if err != nil {
		return err
	}
	project.addEnvFiles()
	addServerlessFile()
	return nil
}

func (p *Project) createRootProjectFolder() error {
	err := os.Mkdir(p.Path, 0777)
	if err != nil {
		return cli.NewExitError("Could not create folder", 1)
	}
	return nil
}

func (p *Project) initProjectGitRepository() error {
	repository, err := git.PlainInit(p.Path, false)
	if err != nil {
		return cli.NewExitError("Could not init git repository", 1)
	}
	p.Repository = repository
	return nil
}

func (p *Project) addSubFolders() error {
	folders := []string{
		"/src/unit_test",
		"/src/integration_test",
		"/conf_git",
	}
	for _, folder := range folders {
		err := os.MkdirAll(path.Join(p.Path, folder), 0777)
		if err != nil {
			return cli.NewExitError("Could not create folder "+folder, 1)
		}
	}
	return nil
}

func (p *Project) addEnvFiles() {
	box := packr.NewBox("../static")
	envTemplate := box.String("template_env")
	for _, env := range ENVS {
		t, _ := template.New("tmp").Parse(envTemplate)
		file, _ := os.Create(path.Join(p.Path, "/conf_git/", env+".env"))
		t.Execute(file, EnvTemplate{EnvName: env})
		file.Close()
	}
}

func CheckExecutableInPath(executable string) bool {
	_, err := exec.LookPath(executable)
	return err == nil
}

func addServerlessFile() {
	config := ServerlessConfig{
		Service: "service",
		Provider: ProviderServerlessConfig{
			Name:       "aws",
			Region:     "eu-west-1",
			Runtime:    "python3.6",
			MemorySize: 128,
		},
	}
	result, _ := yaml.Marshal(&config)
	fmt.Println(string(result))
}
