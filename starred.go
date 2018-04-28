package main

import (
	"context"
	"fmt"
	"github.com/google/go-github/github"
	"os"
	"os/exec"
)

func starredBackupAll(repos []*github.StarredRepository) {
	err := exec.Command("command", "-v", "git").Run()
	if err != nil {
		Error.Println("Failed to backup starred repos. command 'git' not found")
		return
	}

	dir := config.GetString("starred.backupdir")
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.MkdirAll(dir, 0755)
	}

	for _, repo := range repos {
		if config.GetBool("starred.prompt") {
			if prompt(fmt.Sprintf("Backup starred repo %s ?", *repo.Repository.FullName)) {
				reposBackup(repo.Repository, dir, "starred")
			}
		} else {
			reposBackup(repo.Repository, dir, "starred")
		}
	}
}

func starred(ctx context.Context, client *github.Client, username *string) {
	opts := &github.ActivityListStarredOptions{}

	repos, resp, err := client.Activity.ListStarred(ctx, *username, opts)

	if err != nil {
		Error.Printf("Could not read starred repos for user %s\n%s\n", *username, err)
		return
	}
	if resp.StatusCode != 200 {
		Error.Printf("Recieved %d response for starred repos endpoint for user %s\n", resp.StatusCode, *username)
		return
	}
	if len(repos) == 0 {
		Info.Printf("No starred repos for %s", *username)
		return
	}
	starredBackupAll(repos)
}
