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
	"path"
	"strings"
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

type Project struct {
	Path       string
	Repository *git.Repository
	Box        packr.Box
	Serverless *serverless.Config
}

func CreateProject(c *cli.Context) error {
	err := utils.CheckRequiredExecutables(c)
	if err != nil {
		return err
	}
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
	project.preparePythonEnv()
	project.makeFirstCommit()
	return nil
}

func LoadProject() (*Project, error) {
	repository, err := git.PlainOpen(".")
	if err != nil {
		return nil, cli.NewExitError("Could not load git repository", 1)
	}
	storageBox := packr.NewBox("../static")
	serverlessConfig, err := serverless.LoadServerlessConfig("./serverless.yml")
	if err != nil {
		log.Println(err)
		return nil, cli.NewExitError("Could not load serverless config", 1)
	}
	return &Project{
		Path:       ".",
		Repository: repository,
		Box:        storageBox,
		Serverless: serverlessConfig,
	}, nil
}

func (p *Project) createRootProjectFolder() error {
	err := os.Mkdir(p.Path, 0777)
	if err != nil {
		log.Println(err)
		return cli.NewExitError("Could not create folder", 1)
	}
	return nil
}

func (p *Project) initProjectGitRepository() error {
	repository, err := git.PlainInit(p.Path, false)
	if err != nil {
		log.Println(err)
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
	envTemplate := p.Box.String("template_env")
	for _, env := range ENVS {
		t, _ := template.New("tmp").Parse(envTemplate)
		file, _ := os.Create(path.Join(p.Path, "/conf_git/", env+".env"))
		t.Execute(file, EnvTemplate{EnvName: env})
		file.Close()
	}
}

func (p *Project) addServerlessFile() {

	cfg := serverless.NewServerlessConfig("")
	file, _ := os.Create(path.Join(p.Path, "serverless.yml"))
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
	pipfile := p.Box.String("Pipfile")
	file, err := os.Create(path.Join(p.Path, "Pipfile"))
	if err != nil {
		return cli.NewExitError("Could not create Pipfile file", 1)
	}
	defer file.Close()
	file.WriteString(pipfile)

	os.Chdir(p.Path)
	cmd := exec.Command("pipenv", "install", "-d")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	return err
}
