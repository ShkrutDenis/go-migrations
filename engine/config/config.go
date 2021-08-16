package config

import (
	"flag"
	"github.com/joho/godotenv"
	"log"
	"os"
)

var Verbose = false

type Config struct {
	initiated bool

	IsRollback bool
	EnvPath    string
	EnvFile    string
	Verbose    bool

	FirstRun  bool
}

func NewConfig() *Config {
	return &Config{}
}

func NewFromConfig(c Config) *Config {
	Verbose = c.Verbose
	return &Config{
		initiated:  true,
		IsRollback: c.IsRollback,
		EnvPath:    c.EnvPath,
		EnvFile:    c.EnvFile,
		Verbose:    c.Verbose,
	}
}

func (c *Config) InitEnv() {
	if Verbose {
		log.Println("load env file from:", c.GetEnvFullPath())
	}
	err := godotenv.Load(c.GetEnvFullPath())
	if err != nil {
		log.Fatalln("Error loading .env file")
	}
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

func (c *Config) InitByFlags() {
	if c.initiated {
		return
	}
	isRollback := flag.Bool("rollback", false, "Flag for init rollback.")
	envPath := flag.String("env-path", "", "Path to .env file.")
	envFile := flag.String("env-file", ".env", "Env file name.")
	verbose := flag.Bool("verbose", false, "Flag for show more info.")
	flag.Parse()
	c.IsRollback = *isRollback
	c.EnvPath = *envPath
	c.EnvFile = *envFile
	c.Verbose = *verbose
	Verbose = c.Verbose
}
