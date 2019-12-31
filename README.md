[![Go Report Card](https://goreportcard.com/badge/github.com/divbhasin/go-lb)](https://goreportcard.com/report/github.com/divbhasin/go-lb) [![Build Status](https://travis-ci.org/divbhasin/go-lb.svg?branch=master)](https://travis-ci.org/divbhasin/go-lb)

# go-lb: A basic load balancer implemented in Go

go-lb is a load balancer implemented with Go. It uses simple round robin to select backends to handle requests, and performs passive health checks to eliminate inactive backends. It uses http's Reverse Proxy to route requests and fetch responses.

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
  -backends string
    	URLs to backends that need to be load balanced, separated by commas
  -port int
    	Port to serve load balancer on (default 3030)
```

## Built With

* [Go](https://golang.org/) - The language used
* [net/http](https://golang.org/pkg/net/http/) - The package used to proxy requests

## Author

* **Divyam Bhasin**

See also the list of [contributors](https://github.com/your/project/contributors) who participated in this project.

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details

## Acknowledgments

* https://kasvith.github.io/posts/lets-create-a-simple-lb-go/
