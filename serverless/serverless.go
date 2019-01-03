package serverless

import (
	"fmt"
	"github.com/Fall1ngStar/sls-app-builder/function"
	"github.com/imdario/mergo"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
)

type Config struct {
	Service   Service
	Package   Package                      `yaml:",omitempty"`
	Provider  Provider                     `yaml:",omitempty"`
	Functions map[string]function.Function `yaml:",omitempty"`
	Plugins   []string                     `yaml:",omitempty"`
}

type Service struct {
	Name         string
	AwsKmsKeyArn string `yaml:",omitempty"`
}

type Package struct {
	Include []string `yaml:",omitempty"`
	Exclude []string `yaml:",omitempty"`
}

type Provider struct {
	Name               string            `yaml:",omitempty"`
	Region             string            `yaml:",omitempty"`
	Runtime            string            `yaml:",omitempty"`
	MemorySize         int               `yaml:",omitempty"`
	Timeout            int               `yaml:",omitempty"`
	DeploymentBucket   string            `yaml:",omitempty"`
	LogRetentionInDays string            `yaml:",omitempty"`
	Environment        map[string]string `yaml:",omitempty"`
}

type Custom struct {
}

func NewConfig(serviceName string) *Config {
	config := Config{
		Service: Service{
			Name: serviceName,
		},
		Provider: Provider{
			Name:             "aws",
			Region:           "eu-west-1",
			Runtime:          "python3.7",
			MemorySize:       128,
			Timeout:          300,
			DeploymentBucket: serviceName + "-deploys",
		},
	}
	return &config
}

func (cfg *Config) ToYaml() string {
	result, _ := yaml.Marshal(cfg)
	return string(result)
}

func LoadConfig(path string) (*Config, error) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	var config Config
	err = yaml.Unmarshal(content, &config)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return &config, nil
}

func (cfg *Config) UpdateConfigFile(path string) {

	var oldConf, newConf map[interface{}]interface{}
	content, _ := ioutil.ReadFile(path)

	err := yaml.Unmarshal(content, &oldConf)
	if err != nil {
		fmt.Println(err)
	}

	out, _ := yaml.Marshal(cfg)
	err = yaml.Unmarshal(out, &newConf)
	if err != nil {
		fmt.Println(err)
	}

	err = mergo.Merge(&oldConf, newConf)
	if err != nil {
		fmt.Println("mergo:", err)
		return
	}

	result, _ := yaml.Marshal(oldConf)
	err = ioutil.WriteFile(path, result, os.FileMode(os.O_RDWR|os.O_CREATE|os.O_TRUNC))
	if err != nil {
		fmt.Println(err)
	}
}
