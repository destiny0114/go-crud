APP_NAME=go-crud

.PHONY: build
build:
	rm -rf build
	GOARCH=arm64 GOOS=darwin go build -o build/${APP_NAME}-macos cmd/main.go
	GOARCH=amd64 GOOS=linux go build -o build/${APP_NAME}-linux cmd/main.go
	GOARCH=amd64 GOOS=windows go build -o build/${APP_NAME}-windows cmd/main.go