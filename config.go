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
			"backupdir":   os.Getenv("HOME") + "/.ghb/backups/gists",
			"backupregex": "",
			"deleteregex": "",
			"retention":   0,
			"fileonly":    true,
			"prompt":      false,
		},
		"starred": map[string]interface{}{
			"backupdir": os.Getenv("HOME") + "/.ghb/backups/starred",
			"shallow":   true,
			"prompt":    false,
		},
		"repos": map[string]interface{}{
			"backupdir":    os.Getenv("HOME") + "/.ghb/backups/repos",
			"shallow":      true,
			"prompt":       false,
			"owner":        true,
			"collaborator": true,
			"orgmember":    false,
		},
	}
	for key, value := range defaults {
		config.SetDefault(key, value)
	}
}

func configSetup() {
	resp := ""
	dir := os.Getenv("HOME") + "/.ghb"

	fmt.Printf("Full path of directory to save config to [ empty for %s ] : ", dir)

	i, err := fmt.Scanln(&resp)
	if i > 0 && err != nil {
		Error.Printf("Failed to read input for config file path\n%s\n", err)
		shutdown(1)
	}

	if resp == "" {
		resp = dir
	}
	if _, err := os.Stat(resp); os.IsNotExist(err) {
		err = os.MkdirAll(resp, 0755)
		if err != nil {
			Error.Printf("Failed to create directory %s\n%s\n", resp, err)
		}
	}

	err = config.WriteConfigAs(resp + "/config.toml")
	if err != nil {
		Error.Printf("Failed to write config file to %s\n%s", resp, err)
		shutdown(1)
	}
	Info.Printf("Saved default config file to %s/config.toml", resp)

	if resp != (dir) {
		Info.Printf("Make edits as desired then rerun gcb with '-c %s/config.toml' to use\n", resp)
	} else {
		Info.Printf("Make edits as desired then rerun gcb to use\n")
	}
	shutdown(0)
}

func configure(filename string, token string) {
	config = viper.New()
	config.AutomaticEnv()

	configDefaults(config, token)

	if filename == "" {
		if _, err := os.Stat(os.Getenv("HOME") + "/.ghb/config.toml"); os.IsNotExist(err) {
			if prompt("No config file found, would you like to create one?") {
				configSetup()
			}
			return
		}
		config.SetConfigFile(os.Getenv("HOME") + "/.ghb/config.toml")
	} else {
		config.SetConfigFile(filename)
	}
	err := config.ReadInConfig()
	if err != nil {
		Error.Printf("Could not parse config file %s\n%s", filename, err)
		shutdown(1)
	}
}
