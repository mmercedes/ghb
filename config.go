package main

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

func configDefaults(config *viper.Viper, token string) {
	defaults := map[string]interface{}{
		"token": token,
		"enterprise": map[string]string{
			"url": "",
		},
		"gists": map[string]interface{}{
			"backupdir":   os.Getenv("HOME") + "/.ghc/backups/gists",
			"backupregex": "",
			"deleteregex": "",
			"retention":   0,
			"fileonly":    true,
			"prompt":      false,
		},
		"starred": map[string]interface{}{
			"backupdir": os.Getenv("HOME") + "/.ghc/backups/starred",
			"shallow":   true,
			"prompt":    false,
		},
	}
	for key, value := range defaults {
		config.SetDefault(key, value)
	}
}

func configure(filename string, token string) {
	config = viper.New()
	config.AutomaticEnv()

	configDefaults(config, token)

	if filename == "" {
		config.SetConfigName("config")
		config.AddConfigPath(os.Getenv("HOME") + "/.ghc")
	} else {
		config.SetConfigFile(filename)
	}
	err := config.ReadInConfig()
	if err != nil {
		fmt.Printf("Could not parse config file %s\n%s", filename, err)
		shutdown(1)
	}
}
