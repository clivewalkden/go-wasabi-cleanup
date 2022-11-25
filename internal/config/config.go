package config

import (
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"log"
	"os"
	"wasabiCleanup/internal/utils"
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
	pflag.BoolP("verbose", "v", false, "Output additional debug messages")
	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)

	log.Println("Reading config file.")
	viper.SetConfigFile("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(utils.UserHome() + "/.wasabiCleanup/")
	viper.AddConfigPath(".")
	//viper.Debug()

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Println("Config file not found. Please make sure you have a config file in $HOME/.wasabiCleanup/ or ")
			os.Exit(1)
		} else {
			// Config file was found but another error was produced
		}
	}

	config := Config{}
	viper.Unmarshal(&config)

	if viper.GetBool("verbose") {
		log.Println("Loaded config from: ", viper.ConfigFileUsed())
		log.Println("Config: ", viper.AllSettings())
	}

	return config
}
