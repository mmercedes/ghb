package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"golang.org/x/oauth2"

	"github.com/google/go-github/github"
	"github.com/spf13/viper"
)

const (
	// USAGE usage string printed
	USAGE = `
ghb

A tool for performing backups and deletions of Github gists and starred repositories
version: %s
commit: %s
source: https://github.com/mmercedes/ghb


`
)

// overriden by goreleaser
var version = "master"

var (
	commit string

	config *viper.Viper

	// Info - stdout log
	Info *log.Logger
	// Error - stderr log
	Error *log.Logger
)

func init() {
	var (
		configFile string
		token      string

		printversion bool
		debug        bool
		nocolor      bool
	)

	flag.StringVar(&token, "token", os.Getenv("GITHUB_TOKEN"), "Github API token")
	flag.StringVar(&configFile, "config", "", "path to configuration file")
	flag.BoolVar(&printversion, "v", false, "print version")
	flag.BoolVar(&debug, "d", false, "run in debug mode")
	flag.BoolVar(&nocolor, "nc", false, "dont color output")

	flag.Usage = func() {
		fmt.Printf(USAGE, version, commit)
		flag.PrintDefaults()
	}

	flag.Parse()

	if printversion {
		fmt.Printf(USAGE, version, commit)
		shutdown(0)
	}

	logout := ""
	logerr := ""
	if !nocolor {
		logout = "\033[1;32m" // light green
		logerr = "\033[0;31m" // red
	}
	if debug {
		Info = log.New(os.Stdout, logout+"[INFO] ", log.Ldate|log.Ltime|log.Lshortfile)
		Error = log.New(os.Stderr, logerr+"[ERROR] ", log.Ldate|log.Ltime|log.Lshortfile)
	} else {
		Info = log.New(os.Stdout, logout, 0)
		Error = log.New(os.Stderr, logerr, 0)
	}

	// config.go
	configure(configFile, token)

	if config.GetString("token") == "" {
		Error.Println("Github token is required but wasn't set via --token flag, JSON config file,  or found via GITHUB_TOKEN environment variable")
		shutdown(1)
	}

}

func main() {
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigchan
		Info.Printf("Received %s, exiting.\n", sig.String())
		shutdown(0)
	}()

	ctx := context.Background()

	auth := oauth2.NewClient(ctx, oauth2.StaticTokenSource(&oauth2.Token{AccessToken: config.GetString("token")}))

	var client *github.Client
	gheURL := config.GetString("enterprise.url")
	if gheURL == "" {
		client = github.NewClient(auth)
	} else {
		var err error
		client, err = github.NewEnterpriseClient(gheURL, gheURL, auth)
		if err != nil {
			Error.Printf("Could not create github enterprise API client\n%s\n", err)
			shutdown(1)
		}
	}

	user, _, err := client.Users.Get(ctx, "")
	if err != nil {
		Error.Fatalf("%v", err)
	}

	// gists.go
	gists(ctx, client, user.Login)
	// starred.go
	starred(ctx, client, user.Login)
	// repos.go
	repos(ctx, client, user.Login)

	shutdown(0)
}

func shutdown(code int) {
	os.Exit(code)
}

func prompt(msg string) bool {
	var resp string

	fmt.Print("\033[0m" + msg + " [y/N]: ")

	_, err := fmt.Scanln(&resp)
	if err != nil {
		Error.Printf("Failed to read input from user prompt\n%s\n", err)
		return false
	}
	return (strings.ToLower(resp) == "y")
}
