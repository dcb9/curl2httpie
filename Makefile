NAME := curl2httpie
PACKAGE_NAME := github.com/dcb9/curl2httpie

VERSION ?= $(shell git symbolic-ref -q --short HEAD || git describe --tags --exact-match)
GITHUB_SHA ?= $(shell git rev-parse --short HEAD)

VERSION := ${VERSION}
COMMIT := ${GITHUB_SHA}
BUILD_AT := `date`

PLATFORM := linux
BUILD_DIR := build
VAR_SETTING := -X \"$(PACKAGE_NAME)/constant.Version=$(VERSION)\" -X \"$(PACKAGE_NAME)/constant.Commit=$(COMMIT)\" -X \"$(PACKAGE_NAME)/constant.BuildAt=$(BUILD_AT)\"
GOBUILD = go build -ldflags="-s -w $(VAR_SETTING)" -trimpath -o $(BUILD_DIR)

release: clean darwin-amd64.zip linux-amd64.zip freebsd-amd64.zip windows-amd64.zip curl2httpie.js

cloneCurlSrc:
	@echo "\033[0;32mCloning curl source code to local...\033[0m"
	(ls curl_src/docs/cmdline-opts &> /dev/null) || git clone --depth=1 https://github.com/curl/curl.git curl_src
	@echo

generateOptions : cloneCurlSrc
	@echo "\033[0;32mGenerating curl option list...\033[0m"
	rm -rf curl/optionList.go
	go run cmd/generateOptions/main.go -path="./curl_src"
	go fmt curl/optionList.go > /dev/null
	@echo

.PHONY: curl2httpie.js
curl2httpie.js :
	@echo "\033[0;32mBuilding curl2httpie.js ...\033[0m"
	go run -ldflags="$(VAR_SETTING)" ./cmd/fillBuildInfo/main.go
	gopherjs build -m -o public/curl2httpie.js ./cmd/curl2httpie.js
	git checkout constant/

initGithooks:
	git config core.hooksPath .githooks

clean:
	@rm -rf $(BUILD_DIR)
	@rm -f curl2httpie
	@rm -f artifacts
	@mkdir artifacts

test:
	go test ./...

%.zip: %
	@zip -du artifacts/$(NAME)-$@ -j $(BUILD_DIR)/$</*
	@echo "\033[0;32m<<< ---- $(NAME)-$@\033[0m"
	@echo

darwin-amd64:
	@echo "\033[0;32mBuilding $(NAME) for $@...\033[0m"
	mkdir -p $(BUILD_DIR)/$@
	GOARCH=amd64 GOOS=darwin $(GOBUILD)/$@

linux-amd64:
	@echo "\033[0;32mBuilding $(NAME) for $@...\033[0m"
	mkdir -p $(BUILD_DIR)/$@
	GOARCH=amd64 GOOS=linux $(GOBUILD)/$@

freebsd-amd64:
	@echo "\033[0;32mBuilding $(NAME) for $@...\033[0m"
	mkdir -p $(BUILD_DIR)/$@
	GOARCH=amd64 GOOS=freebsd $(GOBUILD)/$@

windows-amd64:
	@echo "\033[0;32mBuilding $(NAME) for $@...\033[0m"
	mkdir -p $(BUILD_DIR)/$@
	GOARCH=amd64 GOOS=windows $(GOBUILD)/$@
