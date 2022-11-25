package utils

import (
	"github.com/spf13/viper"
	"log"
	"os"
)

func UserHome() (homeDir string) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	if viper.GetBool("verbose") {
		log.Println("UserHome Directory: ", homeDir)
	}

	return homeDir
}
