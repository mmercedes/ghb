
GITCOMMIT := $(shell git rev-parse --short HEAD)

default:
	dep ensure
	go build -ldflags "-X main.GitCommit=$(GITCOMMIT)"
