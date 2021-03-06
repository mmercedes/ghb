package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"time"

	"github.com/google/go-github/github"
)

func gistsBackup(gist *github.Gist) {
	if !config.GetBool("gists.fileonly") {
		backupDir := config.GetString("gists.backupdir") + "/" + *gist.ID
		if _, err := os.Stat(backupDir); os.IsNotExist(err) {
			output, err := exec.Command("git", "clone", "-q", *gist.GitPullURL, backupDir).CombinedOutput()
			if err != nil {
				Error.Printf("Failed to clone gist '%s' into '%s'\nClone URL: %s\n%s\n", *gist.HTMLURL, backupDir, *gist.GitPullURL, output)
				return
			}
		} else {
			output, err := exec.Command("git", "-C", backupDir, "pull", "origin", "-q").CombinedOutput()
			if err != nil {
				Error.Printf("Failed to pull remote changes to gist '%s' into '%s'\nPull URL: %s\n%s\n", *gist.HTMLURL, backupDir, *gist.GitPullURL, output)
				return
			}
		}
	} else {
		for _, file := range gist.Files {
			filename := config.GetString("gists.backupdir") + "/" + *gist.ID + "_" + *file.Filename
			output, err := exec.Command("curl", "-s", *file.RawURL, "-o", filename).CombinedOutput()
			if err != nil {
				Error.Printf("Failed to curl gist file %s (%s) into %s\n%s\n", *file.Filename, *file.RawURL, config.GetString("gists.backupdir"), output)
			}
		}
	}
	Info.Printf("Backed up gist '%s' into %s", *gist.HTMLURL, config.GetString("gists.backupdir"))
}

func gistsBackupAll(gists []*github.Gist) {
	command := "curl"
	if !config.GetBool("gists.fileonly") {
		command = "git"
	}

	err := exec.Command("command", "-v", command).Run()
	if err != nil {
		Error.Printf("Failed to backup gists. command '%s' not found\n", command)
		return
	}

	if _, err := os.Stat(config.GetString("gists.backupdir")); os.IsNotExist(err) {
		os.MkdirAll(config.GetString("gists.backupdir"), 0755)
	}

	regex := config.GetString("gists.backupregex")
	for _, gist := range gists {
		if regex != "" {
			for _, file := range gist.Files {
				match, err := regexp.MatchString(regex, *file.Filename)
				if err != nil {
					Error.Printf("Could not match regex %s against string %s\n%s", regex, *file.Filename, err)
					break
				}
				if !match {
					break
				}
				gistsBackup(gist)
			}
		} else {
			gistsBackup(gist)
		}
	}
}

func gistsCanDelete(gist *github.Gist, cutoff time.Time) bool {
	if gist.UpdatedAt.After(cutoff) {
		return false
	}
	regex := config.GetString("gists.deleteregex")
	if regex != "" {
		for _, file := range gist.Files {
			match, err := regexp.MatchString(regex, *file.Filename)
			if err != nil {
				Error.Printf("Could not match regex %s against string %s\n%s", regex, *file.Filename, err)
				return false
			}
			if !match {
				return false
			}
		}
	}
	if config.GetBool("gists.prompt") {
		return prompt(fmt.Sprintf("Delete gist %s ?", *gist.HTMLURL))
	}
	return true
}

func gistsDelete(ctx context.Context, gists []*github.Gist, client *github.Client) {
	if config.GetInt("gists.retention") == 0 {
		return
	}

	cutoff := time.Now().AddDate(0, 0, -config.GetInt("gists.retention"))
	deleted := 0
	for _, gist := range gists {
		if !gistsCanDelete(gist, cutoff) {
			continue
		}
		response, err := client.Gists.Delete(ctx, *gist.ID)
		if err != nil {
			Error.Printf("Failed to delete gist %s\n%s", *gist.HTMLURL, err)
		}
		if response.StatusCode != 204 {
			Error.Printf("Received %d response when attempting to delete gist %s", response.StatusCode, *gist.HTMLURL)
		}
		deleted++
	}
	Info.Printf("Deleted %d gists with no updates after %s", deleted, cutoff.String())
}

func gists(ctx context.Context, client *github.Client, username *string) {
	opts := &github.GistListOptions{Since: time.Time{}}

	gists, resp, err := client.Gists.List(ctx, *username, opts)

	if err != nil {
		Error.Printf("Could not read gists for user %s\n %s\n", *username, err)
		return
	}
	if resp.StatusCode != 200 {
		Error.Printf("Recieved %d response for gists endpoint for user %s\n", resp.StatusCode, *username)
		return
	}
	if len(gists) == 0 {
		Info.Printf("No gists found for %s", *username)
		return
	}
	gistsBackupAll(gists)
	gistsDelete(ctx, gists, client)
}
