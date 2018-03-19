package main

import (
	"context"
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

var (
	token     string
	GitCommit string
	
	version   bool

	Info      *log.Logger
	Error     *log.Logger
)

func init() {
	// TODO: add logfile option
	Info = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)
	Error = log.New(os.Stderr, "[ERROR] ", log.Ldate|log.Ltime|log.Lshortfile)
	
	flag.StringVar(&token, "token", os.Getenv("GITHUB_TOKEN"), "Github API token")

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

	if (token == "") {
		Error.Println("Github token is required but wasn't set via --token flag or found via GITHUB_TOKEN environment variable")
		os.Exit(1)
	}
}

func main() {
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigchan
		Info.Println("Received %s, exiting.", sig.String()) 
		os.Exit(0)
	}()

	ctx := context.Background()

	auth := oauth2.NewClient(ctx, oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token}))
	client := github.NewClient(auth)

	user, _, err := client.Users.Get(ctx, "")
	if err != nil {
		Error.Fatalf("%v", err)
	}

	Info.Printf("%v\n", github.Stringify(user))
	os.Exit(0)
}
