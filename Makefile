GOCMD=go
GOBUILD=$(GOCMD) build
BINARY_NAME=go-lb

build:
		$(GOBUILD) -o $(BINARY_NAME) -v