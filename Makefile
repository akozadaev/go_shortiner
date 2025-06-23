MODULE = $(shell go list -m)
all: build
generate:
	go generate ./...

build: # build a server
	CGO_ENABLED=0 go build -a -o go_shurtiner $(MODULE)/cmd/go_shurtiner

release:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags "-s -w" -o go_shurtiner $(MODULE)/cmd/go_shurtiner/
	zip -9 -r ./go_shurtiner.zip go_shurtiner

lint:
	gofmt -l .

doc:
	go install github.com/swaggo/swag/cmd/swag@latest
	swag init -d ./cmd/go_shurtiner/ -o ./docs --pd

