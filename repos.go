package main

import (
	"context"
	"fmt"
	"github.com/google/go-github/github"
	"os"
	"os/exec"
)

func reposBackup(repo *github.Repository, dir string, confkey string) {
	dir = dir + "/" + *repo.Owner.Login
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.MkdirAll(dir, 0755)
	}

	dir = dir + "/" + *repo.Name
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		cloneURL := *repo.CloneURL
		var output []byte
		var err error

		if config.GetBool(confkey + ".shallow") {
			output, err = exec.Command("git", "clone", "-q", "--depth", "1", cloneURL, dir).CombinedOutput()
		} else {
			output, err = exec.Command("git", "clone", "-q", cloneURL, dir).CombinedOutput()
		}
		if err != nil {
			Error.Printf("Failed to clone repo '%s' into '%s'\n%s\n", cloneURL, dir, output)
			return
		}

	} else {
		if config.GetBool(confkey + ".shallow") {
			output, err := exec.Command("git", "-C", dir, "fetch", "-q", "--depth", "1", "origin").CombinedOutput()
			if err != nil {
				Error.Printf("Failed to fetch remote for clone in %s\n%s\n", dir, output)
				return
			}
			output, err = exec.Command("git", "-C", dir, "reset", "-q", "--hard", "@{upstream}").CombinedOutput()
			if err != nil {
				Error.Printf("Failed to reset to upstream branch for clone in %s\n%s\n", dir, output)
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
	Info.Printf("Backed up repo %s into %s\n", *repo.FullName, dir)
}

func reposBackupAll(repos []*github.Repository) {
	err := exec.Command("command", "-v", "git").Run()
	if err != nil {
		Error.Println("Failed to backup repositories. command 'git' not found")
		return
	}

	dir := config.GetString("repos.backupdir")
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.MkdirAll(dir, 0755)
	}

	for _, repo := range repos {
		if config.GetBool("repos.prompt") {
			if prompt(fmt.Sprintf("Backup repo %s ?", *repo.FullName)) {
				reposBackup(repo, dir, "repos")
			}
		} else {
			reposBackup(repo, dir, "repos")
		}
	}
}

func repos(ctx context.Context, client *github.Client, username *string) {
	repoTypes := ""

	if config.GetBool("repos.owner") {
		repoTypes = "owner"
	}
	if config.GetBool("repos.collaborator") {
		if len(repoTypes) > 0 {
			repoTypes = repoTypes + ","
		}
		repoTypes = repoTypes + "collaborator"
	}
	if config.GetBool("repos.orgmember") {
		if len(repoTypes) > 0 {
			repoTypes = repoTypes + ","
		}
		repoTypes = repoTypes + "organization_member"
	}

	opts := &github.RepositoryListOptions{Affiliation: repoTypes}

	repos, resp, err := client.Repositories.List(ctx, *username, opts)

	if err != nil {
		Error.Printf("Could not read repos for user %s\n%s\n", *username, err)
		return
	}
	if resp.StatusCode != 200 {
		Error.Printf("Recieved %d response for repos endpoint for user %s\n", resp.StatusCode, *username)
		return
	}
	if len(repos) == 0 {
		Info.Printf("No repos for %s", *username)
		return
	}
	reposBackupAll(repos)
}
