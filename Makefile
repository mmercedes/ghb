
GITCOMMIT := $(shell git rev-parse --short HEAD)

.PHONY: install
install:
	go get -u golang.org/x/lint/golint
	go get -u honnef.co/go/tools/cmd/megacheck
	go get -u github.com/golang/dep/cmd/dep
	$(GOPATH)/bin/dep ensure

.PHONY: lint
lint:
	go fmt
	go vet
	$(GOPATH)/bin/golint -set_exit_status
	$(GOPATH)/bin/megacheck -unused.exit-non-zero -simple.exit-non-zero -staticcheck.exit-non-zero

.PHONY: build
build: clean lint
	$(GOPATH)/bin/dep ensure
	go build -ldflags "-X main.commit=$(GITCOMMIT)"

.PHONY: clean
clean:
	@rm -f ghb
