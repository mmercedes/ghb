package main

import (
	"context"
	"os"
	"os/exec"
	"time"

	"github.com/google/go-github/github"
)

func gistsBackup(gist *github.Gist) {
	backupDir := config.BackupDir +"/"+ *gist.ID

	if _, err := os.Stat(backupDir); os.IsNotExist(err) {
		output, err := exec.Command("git", "clone", "-q", *gist.GitPullURL, config.BackupDir+"/"+*gist.ID).CombinedOutput()
		if (err != nil) {
			Error.Printf("Failed to clone gist '%s' into '%s'\nClone URL: %s\n%s\n", *gist.ID, backupDir, *gist.GitPullURL, output)
			return
		}
	} else {
		output, err := exec.Command("git", "-C", backupDir, "pull", "-q").CombinedOutput()
		if (err != nil) {
			Error.Printf("Failed to pull remote changes to gist '%s' into '%s'\nPull URL: %s\n%s\n", *gist.ID, backupDir, *gist.GitPullURL, output)
			return
		}
	}
	Info.Printf("Backed up gist '%s' into %s", *gist.ID, backupDir)
}

func gistsBackupAll(ctx context.Context, client *github.Client, username *string) {
	opts := &github.GistListOptions{Since: time.Time{}}

	gists, response, err := client.Gists.List(ctx, *username, opts)

	if (err != nil) {
		Error.Printf("Could not read gists for user %s\n %s\n", *username, err)
		return
	}
	if (response.StatusCode != 200) {
		Error.Printf("Revied %d response for list gists endpoint for user %s.\n", response.StatusCode, *username)
		return
	}
	if (len(gists) == 0) {
		Info.Printf("No gists found for %s", *username)
		return
	}

	if _, err := os.Stat(config.BackupDir); os.IsNotExist(err) {
		os.MkdirAll(config.BackupDir, 0755)
	}

	for _, gist := range gists {
		gistsBackup(gist)
		//Info.Printf("[%d] %s", i, *gist.HTMLURL)
	}
	return
}

func gists(ctx context.Context, client *github.Client, username *string) {
	err := exec.Command("command", "-v", "git").Run()
	if (err != nil) {
		Error.Println("Failed to backup gists. `git` not found in $PATH")
		return
	}
	gistsBackupAll(ctx, client, username)
}
