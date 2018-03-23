BINARY_NAME=gcp-broker-proxy
BINARY_LINUX=gcp-broker-proxy-linux

all: test build

build:
	go build -o $(BINARY_NAME) -v

test:
	ginkgo .

clean:
	go clean
	rm -rf $(BINARY_NAME)
	rm -rf $(BINARY_LINUX)

run: build
	./$(BINARY_NAME)

build-linux:
				CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o $(BINARY_LINUX) -v

