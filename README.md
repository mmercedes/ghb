# ghb

[![Build Status](https://travis-ci.org/mmercedes/ghb.svg?branch=master)](https://travis-ci.org/mmercedes/ghb)

A simple CLI tool for backing up Github gists and starred respositories

### Features
- works for github.com and enterprise github instances
- backup all gists or those matching a regular expression
- delete gists older than x number of days or matching a regular expression
- clone all of your starred respositories into a backup folder 

### Usage
```
➜ ghb -h

  -config
      path to configuration file (default ~/.ghb/config.toml)
  -d
      run in debug mode
  -h
      show usage
  -nc
      dont color output
  -token
      Github API token (default $GITHUB_TOKEN env var)
  -v
      print version
```

### Configuration

ghb's config file is written in [TOML](https://github.com/toml-lang/toml) and a documented example of the default configruration is provided in [config.toml](https://github.com/mmercedes/ghb/blob/master/config.toml)

### Install

#### Binaries
Linux and OSX binaries can be found under [releases](https://github.com/mmercedes/ghb/releases)

#### Homebrew
```
brew install mmercedes/ghb/ghb
```

#### Go
```
go get github.com/mmercedes/ghb
```
