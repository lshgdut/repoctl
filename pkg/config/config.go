package config

import (
	"fmt"
	"os"

	yaml "gopkg.in/yaml.v3"
)

type Config struct {
	GitlabToken string `yaml:"token"`
	GitlabUrl   string `yaml:"url"`
}

var ConfigFile string = ".repoctl.yaml"

func LoadRepoctlConfig() (*Config, error) {
	yamlConfig := &Config{}
	yamlBytes, err := os.ReadFile(ConfigFile)
	if err != nil {
		return yamlConfig, err
	}

	err = yaml.Unmarshal(yamlBytes, yamlConfig)
	if err != nil {
		return yamlConfig, err
	}

	token := os.Getenv("GITLAB_TOKEN")
	if token == "" {
		token = yamlConfig.GitlabToken
	}

	url := os.Getenv("GITLAB_URL")
	if url == "" {
		url = yamlConfig.GitlabUrl
	}

	if token == "" || url == "" {
		return nil, fmt.Errorf("gitlab token or url not found in config file or environment variables")
	}

	return yamlConfig, nil
}

func (c *Config) Save(dryRun bool) error {
	yamlData := &Config{
		GitlabToken: c.GitlabToken,
		GitlabUrl:   c.GitlabUrl,
	}

	yamlBytes, err := yaml.Marshal(yamlData)
	if err != nil {
		return err
	}

	if dryRun {
		fmt.Println(string(yamlBytes))
		return nil
	}

	err = os.WriteFile(ConfigFile, yamlBytes, 0644)
	if err != nil {
		return err
	}
	return nil
}

func init() {
	config := os.Getenv("REPOCTL_CONFIG")
	if config != "" {
		ConfigFile = config
	}
}
