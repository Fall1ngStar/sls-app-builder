package project

import (
	"github.com/Fall1ngStar/sls-app-builder/serverless"
	"github.com/Fall1ngStar/sls-app-builder/utils"
	"github.com/gobuffalo/packr"
	"github.com/urfave/cli"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var ENVS = []string{
	"DEV",
	"QUAL",
	"PROD",
}

type EnvTemplate struct {
	EnvName string
}

type Project struct {
	Repository  *git.Repository
	Box         packr.Box
	Serverless  *serverless.Config
	ProjectName string
}

func CreateProject(c *cli.Context) error {
	err := utils.CheckRequiredExecutables(c)
	if err != nil {
		return err
	}

	if c.NArg() != 1 {
		return cli.NewExitError("Usage: slapp create <project name>", 1)
	}

	var project Project
	project.ProjectName = c.Args()[0]
	err = project.createRootProjectFolder()
	if err != nil {
		return err
	}

	err = project.initProjectGitRepository()
	if err != nil {
		return err
	}

	project.Box = packr.NewBox("../static")
	err = project.addSubFolders()
	if err != nil {
		return err
	}
	project.addEnvFiles()
	project.addServerlessFile()
	if c.Bool("skip-pipenv") {
		project.preparePythonEnv()
	}
	project.makeFirstCommit()
	return nil
}

func LoadProject() (*Project, error) {
	repository, err := git.PlainOpen(".")
	if err != nil {
		return nil, cli.NewExitError("Could not load git repository", 1)
	}
	storageBox := packr.NewBox("../static")
	serverlessConfig, err := serverless.LoadConfig("./serverless.yml")
	if err != nil {
		log.Println(err)
		return nil, cli.NewExitError("Could not load serverless config", 1)
	}
	return &Project{
		Repository: repository,
		Box:        storageBox,
		Serverless: serverlessConfig,
	}, nil
}

func (p *Project) createRootProjectFolder() error {
	err := os.Mkdir(p.ProjectName, 0777)
	if err != nil {
		log.Println(err)
		return cli.NewExitError("Could not create folder", 1)
	}
	os.Chdir(p.ProjectName)
	return nil
}

func (p *Project) initProjectGitRepository() error {
	repository, err := git.PlainInit(".", false)
	if err != nil {
		log.Println(err)
		return cli.NewExitError("Could not init git repository", 1)
	}
	p.Repository = repository
	return nil
}

func (p *Project) addSubFolders() error {
	folders := []string{
		"src/unit_test",
		"src/integration_test",
		"environ",
	}
	for _, folder := range folders {
		err := os.MkdirAll(folder, 0777)
		if err != nil {
			return cli.NewExitError("Could not create folder "+folder, 1)
		}
	}
	return nil
}

func (p *Project) addEnvFiles() {
	envTemplate, _ := p.Box.FindString("template_env")
	for _, env := range ENVS {
		utils.WriteTemplateToFile(
			envTemplate,
			filepath.Join("environ", env+".yml"),
			EnvTemplate{EnvName: env})
	}
}

func (p *Project) addServerlessFile() {

	cfg := serverless.NewConfig(p.ProjectName)
	file, _ := os.Create("serverless.yml")
	defer file.Close()
	file.WriteString(cfg.ToYaml())
}

func (p *Project) makeFirstCommit() {
	tree, _ := p.Repository.Worktree()
	tree.Add(".")
	tree.Commit("Initial commit", &git.CommitOptions{
		Author: &object.Signature{
			Name: "SLS App Builder CLI",
		},
	})
}

func (p *Project) GetBranchName() string {
	head, _ := p.Repository.Head()
	name := head.Name().Short()
	name = strings.Replace(name, "_", "-", -1)
	result := strings.Join(strings.Split(strings.Title(name), "-")[1:], "")
	return result[:utils.Min(10, len(result))]
}

func (p *Project) preparePythonEnv() error {
	log.Println("Preparing python env")
	pipfile, _ := p.Box.FindString("Pipfile")
	file, err := os.Create("Pipfile")
	if err != nil {
		return cli.NewExitError("Could not create Pipfile file", 1)
	}
	defer file.Close()
	file.WriteString(pipfile)

	cmd := exec.Command("pipenv", "install", "-d")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	return err
}
