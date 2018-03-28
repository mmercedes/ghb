package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/oauth2"
	
	"github.com/google/go-github/github"
)

const (
	USAGE = `ghc
https://github.com/mmercedes/ghc
Git commit: %s
`
)

type Config struct {
	Token     string
	BackupDir string
	LogFile   string
}

var (
	GitCommit string
	
	config    Config

	Info      *log.Logger
	Error     *log.Logger
)


func initConfig(filename string, token string) {
	config = Config{
		Token: token,
		BackupDir: os.Getenv("HOME")+"/ghc_backups",
		LogFile: "",
	}

	if (filename != "") {
		file, _ := os.Open(filename)
		defer file.Close()
		decoder := json.NewDecoder(file)
		err := decoder.Decode(&config)

		if (err != nil) {
			fmt.Printf("Could not parse config file %s\n %s\n", filename, err)
			os.Exit(1)
		} else {
			fmt.Printf("Successfully parsed config file %s. Result:\n %+v\n", filename, config)
		}
	}
}

func init() {
	var configFile string
	var token string
	var version bool
	
	flag.StringVar(&token, "token", os.Getenv("GITHUB_TOKEN"), "Github API token")
	flag.StringVar(&configFile, "config", "", "JSON configuration file full path")
	flag.BoolVar(&version, "v", false, "print version")

	flag.Usage = func() {
		fmt.Printf(USAGE, GitCommit)
		flag.PrintDefaults()
	}

	flag.Parse()

	if (version) {
		fmt.Printf(USAGE, GitCommit)
		os.Exit(0)
	}

	initConfig(configFile, token)

	if (config.LogFile != "") {
		_, err := os.Stat(config.LogFile);
		var logfile *os.File
		if (os.IsNotExist(err)) {
			logfile, err = os.Create(config.LogFile)
		} else {
			logfile, err = os.Open(config.LogFile)
		}
		if (err != nil) {
			fmt.Printf("[ERROR] Could not open logfile %s\n %s", config.LogFile, err);
			os.Exit(1)
		}
		Info = log.New(logfile, "[INFO] ", log.Ldate|log.Ltime|log.Lshortfile)
		Error = log.New(logfile, "[ERROR] ", log.Ldate|log.Ltime|log.Lshortfile)
	} else {
		Info = log.New(os.Stdout, "[INFO] ", log.Ldate|log.Ltime|log.Lshortfile)
		Error = log.New(os.Stderr, "[ERROR] ", log.Ldate|log.Ltime|log.Lshortfile)
	}

	if (config.Token == "") {
		Error.Println("Github token is required but wasn't set via --token flag, JSON config file,  or found via GITHUB_TOKEN environment variable")
		os.Exit(1)
	}
}

func main() {
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigchan
		Info.Printf("Received %s, exiting.\n", sig.String()) 
		os.Exit(0)
	}()

	ctx := context.Background()

	auth := oauth2.NewClient(ctx, oauth2.StaticTokenSource(&oauth2.Token{AccessToken: config.Token}))
	client := github.NewClient(auth)

	user, _, err := client.Users.Get(ctx, "")
	if (err != nil) {
		Error.Fatalf("%v", err)
	}

	// gists.go
	gists(ctx, client, user.Login)

	os.Exit(0)
}

