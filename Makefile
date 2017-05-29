travis-pages :
	go get github.com/gopherjs/gopherjs
	go get ./...
	cd web && gopherjs build -m -o curl2httpie.js && rm main.go .gitignore

release :
	cd cmd/ && goxc -pv="$(v)" -d="$(dest)"

generateOptions :
	go run cmd/generateOptions/main.go -path="$(path)"
	go-bindata -ignore .gitignore -pkg curl -o ./curl/bindata.go data/
