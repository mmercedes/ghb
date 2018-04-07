package main

import (
	"context"
	"fmt"
	"github.com/google/go-github/github"
	"os"
	"os/exec"
)

func starredBackup(repo *github.StarredRepository, dir string) {
	dir = dir + "/" + *repo.Repository.Owner.Login
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.MkdirAll(dir, 0755)
	}

	dir = dir + "/" + *repo.Repository.Name
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		cloneUrl := *repo.Repository.CloneURL
		if config.GetBool("starred.shallow") {
			output, err := exec.Command("git", "clone", "-q", "--depth", "1", cloneUrl, dir).CombinedOutput()
			if err != nil {
				Error.Printf("Failed to shallow clone repo '%s' into '%s'\n%s\n", cloneUrl, dir, output)
				return
			}
		} else {
			output, err := exec.Command("git", "clone", "-q", cloneUrl, dir).CombinedOutput()
			if err != nil {
				Error.Printf("Failed to clone repo '%s' into '%s'\n%s\n", cloneUrl, dir, output)
				return
			}
		}
	} else {
		if config.GetBool("starred.shallow") {
			output, err := exec.Command("git", "-C", dir, "fetch", "-q", "--depth", "1", "origin").CombinedOutput()
			if err != nil {
				Error.Printf("Failed to fetch remote for shallow clone in %s\n%s\n", dir, output)
				return
			}
			output, err = exec.Command("git", "-C", dir, "reset", "-q", "--hard", "@{upstream}").CombinedOutput()
			if err != nil {
				Error.Printf("Failed to reset to upstream branch for shallow clone in %s\n%s\n", dir, output)
				return
			}
		} else {
			output, err := exec.Command("git", "-C", dir, "pull", "-q", "origin").CombinedOutput()
			if err != nil {
				Error.Printf("Failed to pull origin for clone in %s\n%s", dir, output)
				return
			}
		}
	}
	Info.Printf("Backed up starred repo %s into %s\n", *repo.Repository.FullName, dir)
}

func starredBackupAll(repos []*github.StarredRepository) {
	err := exec.Command("command", "-v", "git").Run()
	if err != nil {
		Error.Println("Failed to backup starred repos. command 'git' not found\n")
		return
	}

	dir := config.GetString("starred.backupdir")
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.MkdirAll(dir, 0755)
	}

	for _, repo := range repos {
		if config.GetBool("starred.prompt") {
			if prompt(fmt.Sprintf("Backup starred repo %s ?", *repo.Repository.FullName)) {
				starredBackup(repo, dir)
			}
		} else {
			starredBackup(repo, dir)
		}
	}
}

func starred(ctx context.Context, client *github.Client, username *string) {
	opts := &github.ActivityListStarredOptions{}

	repos, resp, err := client.Activity.ListStarred(ctx, *username, opts)

	if err != nil {
		Error.Printf("Could not read starred repose for user %s\n%s\n", *username, err)
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
