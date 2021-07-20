.PHONY: build clean deploy format

build:
	GOOS=linux GOARCH=amd64 go build -o bin/generatelgtmimage ./cmd/lambda/generateLgtmImage/main.go

clean:
	rm -rf ./bin

deploy: clean build
	npm run deploy

remove:
	npm run remove
