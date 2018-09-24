package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type ServerlessConfig struct {
	Service  map[string]interface{}
	Provider ProviderServerlessConfig `yaml:",omitempty"`
	Custom   CustomServerlessConfig   `yaml:",omitempty"`
}

type ProviderServerlessConfig struct {
	Name             string
	Region           string
	Runtime          string
	MemorySize       int
	Timeout          int
	DeploymentBucket string
	Environment      map[string]string
}

type CustomServerlessConfig struct {
	//Pip string
}

func NewServerlessConfig() *ServerlessConfig {
	config := ServerlessConfig{
		Service: map[string]interface{}{},
		Provider: ProviderServerlessConfig{
			Name:       "aws",
			Region:     "eu-west-1",
			Runtime:    "python3.6",
			MemorySize: 128,
			Timeout:    300,
			Environment: map[string]string{
				//"ACCOUNT_ID": "${env:ACCOUNT_ID}",
			},
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
		return nil, err
	}
	var config ServerlessConfig
	err = yaml.Unmarshal(content, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}