package config

import (
	"fmt"
	"github.com/spf13/viper"
	"log"
	"os"
)

type Config struct {
	Buckets    map[string]int `yaml:"buckets"`
	Connection S3Connection   `yaml:"connection"`
}

type S3Connection struct {
	Url     string `yaml:"url"`
	Region  string `yaml:"region"`
	Profile string `yaml:"profile"`
}

func InitConfig() Config {
	fmt.Fprintln(os.Stdout, "Reading config file")
	// Find home directory.
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	viper.SetConfigFile("wasabi-cleanup.yml")
	viper.AddConfigPath(".")
	viper.AddConfigPath(home)
	viper.AutomaticEnv()

	// If a config file is found, read it in.
	err = viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	config := Config{}
	viper.Unmarshal(&config)

	//fmt.Println(config)

	return config
}
