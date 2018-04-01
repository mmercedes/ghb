package main

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
)

type Config struct {
	Token         string
	BackupDir     string
	EnterpriseUrl string
	
	DeleteAfter   int

	FullBackup    bool
}

func configure(filename string, token string) {
	config = Config{}
	defaults := Config{
		Token: token,
		BackupDir: os.Getenv("HOME")+"/.ghc/backups",
		EnterpriseUrl: "",
		DeleteAfter: 0,
		FullBackup: false,
	}

	// check for config file in default location
	if (filename == "") {
		filename = os.Getenv("HOME")+"/.ghc/conf.json";
		if _, err := os.Stat(filename); os.IsNotExist(err) {
			fmt.Printf("No config file found in %s using defaults\n", filename)
			config = defaults
			return
		} 
	}

	file, _ := os.Open(filename)
	defer file.Close()
	decoder := json.NewDecoder(file)
	err := decoder.Decode(&config)
	
	if (err != nil) {
		fmt.Printf("Could not parse config file %s\n %s\n", filename, err)
		shutdown(1)
	}

	// if the parsed config has any empty string options, set them to the defaults
	rconfig := reflect.ValueOf(&config).Elem()
	rdefault := reflect.ValueOf(&defaults).Elem()
	
	for i := 0; i < rconfig.NumField(); i++ {
		field := rconfig.Field(i)

		if (field.Type() != reflect.TypeOf("")) {
			continue
		}
		if (field.Interface().(string) != "") {
			continue
		}

		defstr := rdefault.Field(i).Interface().(string)
		field.SetString(defstr)
	}

	fmt.Printf("Successfully parsed config file %s. Result:\n %+v\n", filename, config)
}

