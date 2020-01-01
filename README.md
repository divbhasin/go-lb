[![Go Report Card](https://goreportcard.com/badge/github.com/divbhasin/go-lb)](https://goreportcard.com/report/github.com/divbhasin/go-lb) [![Build Status](https://travis-ci.org/divbhasin/go-lb.svg?branch=master)](https://travis-ci.org/divbhasin/go-lb) [![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

# go-lb: A basic load balancer implemented in Go

go-lb is a load balancer implemented with Go. It uses simple round robin to select backends to handle requests, and performs passive health checks to eliminate inactive backends. It uses http's Reverse Proxy to route requests and fetch responses. Please note that this project is still a work-in-progress.

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes.

### Prerequisites

Go language
Built-in Go packages

Note: Make sure your GOBIN and GOPATH environment variables are set up appropriately. For example:

```
export GOPATH=$HOME/go
export GOBIN=$GOPATH/bin
```

### Installing

To install, first run
```
go get github.com/divbhasin/go-lb
```

### Usage
```
Usage of go-lb:
  -config string
    	The name of the config file to use (default: config.yml)
```

A sample config file:
```
port: 3030
backends:
  - localhost:3031
  - localhost:3031
```

## Built With

* [Go](https://golang.org/) - The language used
* [net/http](https://golang.org/pkg/net/http/) - The package used to proxy requests

## Author

* **Divyam Bhasin**

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details

## Acknowledgments

* https://kasvith.github.io/posts/lets-create-a-simple-lb-go/
