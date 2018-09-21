package project

import (
	"github.com/urfave/cli"
	"os"
	"gopkg.in/libgit2/git2go.v26"
	"os/exec"
)

func createRootProjectFolder(path string) error {
	err := os.Mkdir(path, 0777)
	if err != nil {
		return cli.NewExitError("Could not create folder", 1)
	}
	return nil
}

func initProjectGitRepository(path string) (*git.Repository, error) {
	repository, err := git.InitRepository(path, false)
	if err != nil {
		return nil, cli.NewExitError("Could not init git repository", 1)
	}
	return repository, nil
}

func CreateProject(c *cli.Context) error {
	var path string
	if c.NArg() == 0 {
		path = "."
	} else if c.NArg() == 1 {
		path = c.Args()[0]
		err := createRootProjectFolder(path)
		if err != nil {
			return err
		}
	} else {
		return cli.NewExitError("Wrong argument number provided", 1)
	}

	_, err := initProjectGitRepository(path)
	if err != nil {
		return err
	}
	return nil
}

func CheckServerlessRequirements() bool {
	_, err := exec.LookPath("serverless")
	return err == nil
}
