package config

import (
	"os"

	"github.com/go-yaml/yaml"
)

type Config struct {
	// Config struct metadata
	
	Environment string

	// Config fields
	
	BallastSize int64 //`yaml:"ballastSize"` // This annotation should be unnecessary

	// Config collections
	
	Database Database
}
type Database struct {
	DriverName       string
	DataSourceName   string
	KeepAliveSeconds int64
}

func (c *Config) Load(env string) error {
	if yamlFile, err := os.ReadFile("./config/" + env + ".yaml"); err != nil {
		return err
	} else if err := yaml.Unmarshal(yamlFile, c); err != nil {
		return err
	}

	c.Environment = env
	return nil
}
