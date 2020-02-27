SHELL := /bin/bash

VERSION := $(CI_COMMIT_TAG)
GITCOMMIT := `git rev-parse HEAD`
BRANCH := $(CI_COMMIT_BRANCH)
BUILDDATE := `date +%Y-%m-%d`
BUILDUSER := `whoami`

LDFLAGSSTRING :=-X gitlab.com/zcash/zcashd_exporter/version.Version=$(VERSION)
LDFLAGSSTRING +=-X gitlab.com/zcash/zcashd_exporter/version.GitCommit=$(GITCOMMIT)
LDFLAGSSTRING +=-X gitlab.com/zcash/zcashd_exporter/version.Branch=$(BRANCH)
LDFLAGSSTRING +=-X gitlab.com/zcash/zcashd_exporter/version.BuildDate=$(BUILDDATE)
LDFLAGSSTRING +=-X gitlab.com/zcash/zcashd_exporter/version.BuildUser=$(BUILDUSER)

LDFLAGS :=-ldflags "$(LDFLAGSSTRING)"

.PHONY: all build

all: build

# Build binary
build:
	go build $(LDFLAGS) 