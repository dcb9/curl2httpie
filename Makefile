travis-pages :
	go get github.com/gopherjs/gopherjs
	go get ./...
	go test ./...
	cd web && gopherjs build -m -o curl2httpie.js && rm main.go .gitignore

# $ dest=~/Downloads/curl2httpie/ v=v1.x make release
release :
	cd cmd/curl2httpie/ && goxc -pv="$(v)" -d="$(dest)"

generateOptions :
	go run cmd/generateOptions/main.go -path="$(path)"
	go-bindata -ignore .gitignore -pkg curl -o ./curl/bindata.go data/

initGithooks:
	git config core.hooksPath .githooks