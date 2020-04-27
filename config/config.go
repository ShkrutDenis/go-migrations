package config

import "os"

var cfg *Config

type Config struct {
	IsRollback bool
	EnvPath    string
	EnvFile    string
	Verbose    bool

	LastBatch int
	FirstRun  bool
}

func (c *Config) GetEnvFullPath() string {
	if c.EnvPath == "" {
		return c.EnvFile
	}
	_, err := os.Stat(c.EnvPath + "/" + c.EnvFile)
	if os.IsNotExist(err) {
		return c.EnvPath + "\\" + c.EnvFile
	}
	return c.EnvPath + "/" + c.EnvFile
}

func GetConfig() *Config {
	if cfg == nil {
		cfg = &Config{}
	}
	return cfg
}
