package main

import (
	"context"
	"os"
	"os/exec"
	"time"

	"github.com/google/go-github/github"
)

func gistsBackup(gist *github.Gist) {
	if (config.FullBackup) {
		backupDir := config.BackupDir +"/"+ *gist.ID
		if _, err := os.Stat(backupDir); os.IsNotExist(err) {
			output, err := exec.Command("git", "clone", "-q", *gist.GitPullURL, backupDir).CombinedOutput()
			if (err != nil) {
				Error.Printf("Failed to clone gist '%s' into '%s'\nClone URL: %s\n%s\n", *gist.HTMLURL, backupDir, *gist.GitPullURL, output)
				return
			}
		} else {
			output, err := exec.Command("git", "-C", backupDir, "pull", "-q").CombinedOutput()
			if (err != nil) {
				Error.Printf("Failed to pull remote changes to gist '%s' into '%s'\nPull URL: %s\n%s\n", *gist.HTMLURL, backupDir, *gist.GitPullURL, output)
				return
			}
		}
	} else {
		for _, file := range gist.Files {
			filename := config.BackupDir + "/" + *gist.ID + "_" + *file.Filename
			output, err := exec.Command("curl", "-s", *file.RawURL, "-o", filename).CombinedOutput()
			if (err != nil) {
				Error.Printf("Failed to curl gist file %s (%s) into %s\n%s\n", *file.Filename, *file.RawURL, config.BackupDir, output)
			}
		}				
	}
	Info.Printf("Backed up gist '%s' into %s", *gist.HTMLURL, config.BackupDir)
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
	}
	return
}

func gists(ctx context.Context, client *github.Client, username *string) {
	command := "curl"
	if (config.FullBackup) {
		command = "git"
	}

	err := exec.Command("command", "-v", command).Run()
	if (err != nil) {
		Error.Printf("Failed to backup gists. command '%s' not found\n", command)
		return
	}
	gistsBackupAll(ctx, client, username)
}
