# Define the binary name
BINARY_NAME=bin/juno

# Define the build command
build:
	go build -o $(BINARY_NAME) *.go

install:
	export GOBIN=/usr/local/bin/
	go install ./...