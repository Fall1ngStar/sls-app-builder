package config

import (
	"github.com/imdario/mergo"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
)

type ServerlessConfig struct {
	//Service  map[string]interface{}
	Provider Provider `yaml:",omitempty"`
	//Custom   Custom   `yaml:",omitempty"`
}

type Service struct {
	Name         string
	AwsKmsKeyArn string
}

type Package struct {
	Include []string
	Exclude []string
}

type Provider struct {
	Name               string
	Region             string
	Runtime            string
	MemorySize         int
	Timeout            int
	DeploymentBucket   string
	LogRetentionInDays string
	Environment        map[string]string
}

type Custom struct {
}

func NewServerlessConfig() *ServerlessConfig {
	config := ServerlessConfig{
		//Service: map[string]interface{}{},
		Provider: Provider{
			Name:       "aws",
			Region:     "eu-west-1",
			Runtime:    "python3.6",
			MemorySize: 128,
			Timeout:    300,
		},
	}
	return &config
}

func (cfg *ServerlessConfig) ToYaml() string {
	result, _ := yaml.Marshal(cfg)
	return string(result)
}

func LoadServerlessConfig(path string) (*ServerlessConfig, error) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	var config ServerlessConfig
	err = yaml.Unmarshal(content, &config)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return &config, nil
}

func (cfg *ServerlessConfig) UpdateConfigFile(path string) {
	var oldConf, newConf yaml.MapSlice
	content, _ := ioutil.ReadFile(path)
	yaml.Unmarshal(content, &oldConf)

	out, _ := yaml.Marshal(cfg)
	yaml.Unmarshal(out, &newConf)

	mergo.Merge(&oldConf, newConf)
	result, _ := yaml.Marshal(oldConf)
	ioutil.WriteFile("serverless2.yml", result, os.FileMode(os.O_RDWR|os.O_CREATE|os.O_TRUNC))
}
