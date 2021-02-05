NAME := curl2httpie
PACKAGE_NAME := github.com/dcb9/curl2httpie
VERSION := `git describe --dirty`
COMMIT := `git rev-parse HEAD`

PLATFORM := linux
BUILD_DIR := build
VAR_SETTING := -X $(PACKAGE_NAME)/constant.Version=$(VERSION) -X $(PACKAGE_NAME)/constant.Commit=$(COMMIT)
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

curl2httpie.js.deps : generateOptions
	@echo "\033[0;32mInstalling gopherjs and it's deps...\033[0m"
	cd ../ \
	&& which go1.12.16 \
	|| (go get golang.org/dl/go1.12.16 && go1.12.16 download)
	which gopherjs || go get github.com/gopherjs/gopherjs
	@echo

.PHONY: curl2httpie.js
curl2httpie.js : curl2httpie.js.deps
	@echo "\033[0;32mBuilding curl2httpie.js ...\033[0m"
	GOPHERJS_GOROOT="/Users/bob/sdk/go1.12.16" gopherjs build -m -o docs/curl2httpie.js ./cmd/curl2httpie.js

initGithooks:
	git config core.hooksPath .githooks

clean:
	@rm -rf $(BUILD_DIR)
	@rm -f curl2httpie
	@rm -f curl2httpie-*.zip

test:
	go test ./...

%.zip: %
	@zip -du $(NAME)-$@ -j $(BUILD_DIR)/$</*
	@echo "\033[0;32m<<< ---- $(NAME)-$@\033[0m"
	@echo

darwin-amd64: generateOptions
	@echo "\033[0;32mBuilding $(NAME) for $@...\033[0m"
	mkdir -p $(BUILD_DIR)/$@
	GOARCH=amd64 GOOS=darwin $(GOBUILD)/$@

linux-amd64: generateOptions
	@echo "\033[0;32mBuilding $(NAME) for $@...\033[0m"
	mkdir -p $(BUILD_DIR)/$@
	GOARCH=amd64 GOOS=linux $(GOBUILD)/$@

freebsd-amd64: generateOptions
	@echo "\033[0;32mBuilding $(NAME) for $@...\033[0m"
	mkdir -p $(BUILD_DIR)/$@
	GOARCH=amd64 GOOS=freebsd $(GOBUILD)/$@

windows-amd64: generateOptions
	@echo "\033[0;32mBuilding $(NAME) for $@...\033[0m"
	mkdir -p $(BUILD_DIR)/$@
	GOARCH=amd64 GOOS=windows $(GOBUILD)/$@
