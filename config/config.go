package config

import (
	"log"
	"os"
	"strconv"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Wss      string `yaml:"wss"`
	Parallel bool   `yaml:"parallel"`
}

var Conf *Config

func LoadConfig() *Config {
	if Conf != nil {
		return Conf
	}
	Conf = &Config{}
	loadFile()
	loadEnv()
	return Conf
}

/*
** When more variables are added along the way, this method should be refactored.
** Either use a third-party package like viper, envconfig etc. or store the
** variable names in a slice and look for each of them in "os.Environ()".
*/
func loadEnv() {
	parallel := os.Getenv("PARALLEL")
	wss := os.Getenv("WSS")
	if parallel != "" {
		parBool, err := strconv.ParseBool(parallel)
		ifErrExit(err)
		Conf.Parallel = parBool
	}
	if wss != "" {
		Conf.Wss = wss
	}
}

func loadFile() {
	f, err := os.Open("config/config.yml")
	if err != nil {
		log.Printf("Error loading config file, looking for env vars...\n\n")
		return
	}
	defer f.Close()
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&Conf)
	ifErrExit(err)
}

func ifErrExit(err error) {
	if err != nil {
		log.Printf("Error with config file/env vars:\n'%s'\n", err)
		os.Exit(1)
	}
}
