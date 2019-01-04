package project

import (
	"fmt"
	"github.com/Fall1ngStar/sls-app-builder/serverless"
	"github.com/Fall1ngStar/sls-app-builder/utils"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/endpoints"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gobuffalo/packr"
	"github.com/urfave/cli"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"io/ioutil"
	"log"
	"os"
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

	err = project.addServerlessFile()
	if err != nil {
		return err
	}

	err = project.preparePythonEnv()
	if err != nil {
		return err
	}

	err = project.addServerlessPlugins()
	if err != nil {
		return err
	}

	err = project.createDeployBucket()
	if err != nil {
		return err
	}

	err = project.makeFirstCommit()
	if err != nil {
		return err
	}

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

func (p *Project) addServerlessFile() error {

	cfg := serverless.NewConfig(p.ProjectName)
	p.Serverless = cfg
	file, _ := os.Create("serverless.yml")
	defer file.Close()
	_, err := file.WriteString(cfg.ToYaml())
	if err != nil {
		return err
	}
	return nil
}

func (p *Project) makeFirstCommit() error {
	tree, _ := p.Repository.Worktree()
	_, err := tree.Add(".")
	if err != nil {
		return err
	}
	_, err = tree.Commit("Initial commit", &git.CommitOptions{
		Author: &object.Signature{
			Name: "SLS App Builder CLI",
		},
	})
	if err != nil {
		return err
	}
	return nil
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
	_, err = file.WriteString(pipfile)
	if err != nil {
		return err
	}

	return utils.RunWithStdout("pipenv", "install", "-d")
}

func (p *Project) addServerlessPlugins() error {
	packageJson, err := p.Box.FindString("package.json")
	if err != nil {
		fmt.Println(err)
		return err
	}
	err = ioutil.WriteFile("package.json", []byte(packageJson), 0755)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return utils.RunWithStdout("npm", "install")
}

func (p *Project) createDeployBucket() error {
	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		fmt.Println(err)
		return err
	}
	cfg.Region = endpoints.EuWest1RegionID
	s3Client := s3.New(cfg)
	bucketName := aws.String(p.Serverless.Service.Name + "-deploys")
	req := s3Client.CreateBucketRequest(&s3.CreateBucketInput{
		Bucket: bucketName,
	})
	_, err = req.Send()
	// Excluding errors of bucket already existing
	if err != nil && !strings.Contains(err.Error(), "already") {
		return err
	}
	return nil
}
