SHELL := /bin/bash

#===============================================================================
#  release information
#===============================================================================
PACKAGE := $(shell go list)
BINARY := $(notdir $(PACKAGE))

TOOL_DIR := _tool
RELEASE_DIR := _release
PKG_DEST_DIR := $(RELEASE_DIR)/.pkg

ALL_OS := darwin linux windows
ALL_ARCH := 386 amd64

LATEST_LOCAL_DEVEL_BRANCH := $(subst * ,,$(shell git branch --sort='-committerdate' \
	| grep --invert-match master                                                    \
	| head --lines=1))
NEW_TAG := $(shell echo "$(LATEST_LOCAL_DEVEL_BRANCH)" \
	| grep --only-matching -E '[0-9]+\.[0-9]+\.[0-9]+')

#===============================================================================
#  version information embedding
#===============================================================================
# バージョンタグは `git tag -a 'x.y.z'` と注釈付きタグであることが前提。
VERSION := $(shell git describe --always --dirty 2>/dev/null || echo 'no git tag')
VERSION_PACKAGE := github.com/yuta-masano/ghb/cmd
BUILD_REVISION := $(shell git rev-parse --short HEAD)
BUILD_WITH := $(shell go version)
STATIC_FLAGS := -a -tags netgo -installsuffix netgo
LD_FLAGS := -s -w -X '$(VERSION_PACKAGE).buildVersion=$(VERSION)' \
	-X '$(VERSION_PACKAGE).buildRevision=$(BUILD_REVISION)'       \
	-X '$(VERSION_PACKAGE).buildWith=$(BUILD_WITH)'               \
	-extldflags -static

#===============================================================================
#  lint options
#===============================================================================
GOMETALINTER_OPTS := --enable-all --vendored-linters --deadline=60s \
	--dupl-threshold=75 --line-length=120
GOMETALINTER_EXCLUDE_REGEX := gas

#===============================================================================
#  targets
#    `make [help]` shows tasks what you should execute.
#    The other are helper targets.
#===============================================================================
.DEFAULT_GOAL := help

# [Add a help target to a Makefile that will allow all targets to be self documenting]
# https://gist.github.com/prwhite/8168133
.PHONY: help
help: ## show help
	@echo 'USAGE: make [target]'
	@echo
	@echo 'TARGETS:'
	@grep -E '^[-_: a-zA-Z0-9]+##' $(MAKEFILE_LIST) \
		| sed 's/:[-_ a-zA-Z0-9]\+/:/'              \
		| column -t -s ':#'

# install development tools
.PHONY: setup
setup:
	type -a glide &>/dev/null || curl https://glide.sh/get | sh
	go get -v -u github.com/alecthomas/gometalinter
	go get -v -u github.com/tcnksm/ghr
	gometalinter --install
	cp -a $(TOOL_DIR)/git_hooks/* .git/hooks/

.PHONY: deps-install
deps-install: setup ## install vendor packages based on glide.lock or glide.yaml
	glide install --strip-vendor

.PHONY: install
install:
	test -e glide.yaml \
		&& go install $(shell sed --quiet 's/^- package: /.\/vendor\//p' glide.yaml) || :
	CGO_ENABLED=0 go install $(subst -a ,,$(STATIC_FLAGS)) -ldflags "$(LD_FLAGS)"

.PHONY: lint
lint: install ## lint go sources and check whether only LICENSE file has copyright sentence
	gometalinter $(GOMETALINTER_OPTS)                                                  \
		$(if $(GOMETALINTER_EXCLUDE_REGEX), --exclude='$(GOMETALINTER_EXCLUDE_REGEX)') \
		$(shell glide novendor)
	$(TOOL_DIR)/copyright_check.sh

.PHONY: push-release-tag
push-release-tag: lint test ## update CHANGELOG and push all of the your development works
	$(TOOL_DIR)/add_changelog.sh "$(NEW_TAG)"
	git checkout master
	git merge --ff "$(LATEST_LOCAL_DEVEL_BRANCH)"
	git push
	$(TOOL_DIR)/add_release_tag.sh "$(NEW_TAG)"

.PHONY: test
test: ## go test
	go test -v -cover $(shell glide novendor)

.PHONY: all-build
all-build: lint test
	$(TOOL_DIR)/build_static_bins.sh "$(ALL_OS)" "$(ALL_ARCH)" \
		"$(STATIC_FLAGS)" "$(LD_FLAGS)"                        \
		"$(PKG_DEST_DIR)" "$(BINARY)"

.PHONY: all-archive
all-archive:
	$(TOOL_DIR)/archive.sh "$(ALL_OS)" "$(ALL_ARCH)" "$(PKG_DEST_DIR)"

.PHONY: release
release: all-build all-archive ## build binaries for all platforms and upload them to GitHub
	ghr "$(VERSION)" "$(RELEASE_DIR)"

.PHONY: clean
clean: ## uninstall the binary and remove $(RELEASE_DIR) directory
	go clean -i .
	rm -rf $(RELEASE_DIR)
