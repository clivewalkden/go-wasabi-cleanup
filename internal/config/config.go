package config

import (
	"fmt"
	"github.com/spf13/cobra"
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

var cfgFile string

func InitConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		viper.AddConfigPath(home + "/.wasabiCleanup/")
		viper.AddConfigPath(".")
		viper.SetConfigFile("config")
		viper.SetConfigType("yaml")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}

	if viper.GetBool("verbose") {
		log.Println("Loaded config from: ", viper.ConfigFileUsed())
		log.Println("Config: ", viper.AllSettings())
	}
}

func AppConfig() Config {
	config := Config{}
	viper.Unmarshal(&config)

	return config
}
