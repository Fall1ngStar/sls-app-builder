package serverless

import (
	"github.com/imdario/mergo"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
)

type Config struct {
	Service   Service
	Package   Package             `yaml:",omitempty"`
	Provider  Provider            `yaml:",omitempty"`
	Functions map[string]Function `yaml:",omitempty"`
	Plugins   []string            `yaml:",omitempty"`
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

type Function struct {
	Handler string
	Events  []map[string]interface{} `yaml:",omitempty"`
}

type Custom struct {
}

func NewServerlessConfig(serviceName string) *Config {
	config := Config{
		Service: Service{
			Name: serviceName,
		},
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

func (cfg *Config) ToYaml() string {
	result, _ := yaml.Marshal(cfg)
	return string(result)
}

func LoadServerlessConfig(path string) (*Config, error) {
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
	var oldConf, newConf yaml.MapSlice
	content, _ := ioutil.ReadFile(path)
	yaml.Unmarshal(content, &oldConf)

	out, _ := yaml.Marshal(cfg)
	yaml.Unmarshal(out, &newConf)

	mergo.Merge(&oldConf, newConf)
	result, _ := yaml.Marshal(oldConf)
	ioutil.WriteFile(path, result, os.FileMode(os.O_RDWR|os.O_CREATE|os.O_TRUNC))
}
