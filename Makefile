SHELL := /bin/bash

VERSION := ${VERSION}
GITCOMMIT := `git rev-parse HEAD`
BUILDDATE := `date +%Y-%m-%d`
BUILDUSER := `whoami`

LDFLAGSSTRING :=-X gitlab.com/zcash/zcashd_exporter/version.Version=$(VERSION)
LDFLAGSSTRING +=-X gitlab.com/zcash/zcashd_exporter/version.GitCommit=$(GITCOMMIT)
LDFLAGSSTRING +=-X gitlab.com/zcash/zcashd_exporter/version.BuildDate=$(BUILDDATE)
LDFLAGSSTRING +=-X gitlab.com/zcash/zcashd_exporter/version.BuildUser=$(BUILDUSER)

LDFLAGS :=-ldflags "$(LDFLAGSSTRING)"

.PHONY: all build

all: build

# Build binary
build:
	CGO_ENABLED=0 go build $(LDFLAGS) 