# Fetch exact tag (if present)
TAG := $(shell git describe --tags --exact-match 2>/dev/null)
ifeq ($(TAG),)
# Ok, no exact tag, fetch the latest
TAG := $(shell git describe --abbrev=0 --tags 2>/dev/null)
ifeq ($(TAG),)
# Also no latest tag, default to v0.0.0
TAG := "v0.0.0"
endif
# Append the commit-hash when the tag is not exact
TAG := $(shell echo "$(TAG)+git.$(shell git log -n 1 --pretty="%h")")
endif

TARGETS := $(shell go list ./... | grep -v vendor)

all:
	$(MAKE) get
	$(MAKE) build
	$(MAKE) test

get:
	@for target in $(TARGETS); do go get -t -v $$target; done

build:
	go build -ldflags "-X main.buildDate=`date -u +%Y-%m-%d:%H:%M:%S` -X main.buildVersion=$(TAG)"

test:
	@echo "" > coverage.txt
	@for d in $(TARGETS); do \
    	go test -race -coverprofile=profile.out -covermode=atomic $$d ; \
	    if [ -f profile.out ]; then cat profile.out >> coverage.txt; rm -f profile.out; fi ; \
	done