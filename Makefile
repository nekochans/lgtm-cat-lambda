.PHONY: build clean deploy format generate

build:
	GOOS=linux GOARCH=amd64 go build -o bin/generatelgtmimage ./cmd/lambda/generatelgtmimage/main.go
	GOOS=linux GOARCH=amd64 go build -o bin/storelgtmimage ./cmd/lambda/storelgtmimage/main.go

generate:
	docker run --rm -v `pwd`:/src -w /src kjconroy/sqlc generate

clean:
	rm -rf ./bin

deploy: clean build
	npm run deploy

remove:
	npm run remove
