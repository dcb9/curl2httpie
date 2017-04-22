travis-pages :
	go get github.com/gopherjs/gopherjs
	go get ./...
	cd web && gopherjs build -m -o curl2httpie.js && rm main.go .gitignore

release :
	cd cmd/ && goxc -pv="$(v)" -d="$(dest)"
