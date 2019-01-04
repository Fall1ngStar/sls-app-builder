package commands

import (
	"github.com/Fall1ngStar/sls-app-builder/project"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/endpoints"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/urfave/cli"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
)

func DeployProject(c *cli.Context) error {
	p, err := project.LoadProject()
	if err != nil {
		return err
	}

	gitlab := c.Bool("gitlab")
	stage := c.String("stage")
	branch := p.GetBranchName()

	err = loadEnvFile(stage, gitlab)
	if err != nil {
		return err
	}
	if branch == "" {
		branch = stage
	}

	err = getDepLayerVersion(p, branch)
	if err != nil {
	    return err
	}

	cmd := exec.Command("serverless", "deploy", "--stage", branch)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
	    return err
	}
	return nil
}

func loadEnvFile(stage string, gitlab bool) error {
	content, err := ioutil.ReadFile(filepath.Join("environ", stage+".yml"))
	if err != nil {
		return err
	}

	contentWithEnv := os.ExpandEnv(string(content))
	var variables map[string]string
	err = yaml.Unmarshal([]byte(contentWithEnv), &variables)
	if err != nil {
		return err
	}

	for key, value := range variables {
		_ = os.Setenv(key, value)
	}

	if !gitlab {
		_ = os.Unsetenv("AWS_ACCESS_KEY_ID")
		_ = os.Unsetenv("AWS_SECRET_ACCESS_KEY")
	}
	return nil
}

func getDepLayerVersion(p *project.Project, env string) error {
	layerName := p.Serverless.Service.Name + "-deps-" + env
	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
	    return err
	}
	cfg.Region = endpoints.EuWest1RegionID
	lambdaClient := lambda.New(cfg)

	resp, err := lambdaClient.ListLayerVersionsRequest(&lambda.ListLayerVersionsInput{
		LayerName:aws.String(layerName),
	}).Send()
	if err != nil {
	    return err
	}
	if len(resp.LayerVersions) == 0 {
		return cli.NewExitError("Dependency layer is not deployed, please use \"slapp layers\"", 1)
	}

	sort.Slice(resp.LayerVersions, func(i, j int) bool {
		return *resp.LayerVersions[i].Version < *resp.LayerVersions[j].Version
	})
	err = os.Setenv("LAYER_VERSION_ARN", *resp.LayerVersions[0].LayerVersionArn)
	if err != nil {
	    return err
	}
	return nil
}
