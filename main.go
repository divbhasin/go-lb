package main

import (
	"context"
	"flag"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"
)

const (
	attempts int = iota
	retry
)

type config struct {
	Port     int64    `yaml:"port"`
	Backends []string `yaml:"backends"`
}

var s ServerPool

func lb(w http.ResponseWriter, r *http.Request) {
	attempts := GetAttemptsFromContext(r)
	if attempts > 3 {
		log.Printf("%s(%s) Max attempts reached, terminating\n", r.RemoteAddr, r.URL.Path)
		http.Error(w, "Service not available", http.StatusServiceUnavailable)
		return
	}

	peer := s.GetNextPeer()
	if peer != nil {
		peer.ReverseProxy.ServeHTTP(w, r)
		return
	}
	http.Error(w, "Service not available", http.StatusServiceUnavailable)
}

// GetRetryFromContext : takes a request context and deduces the number of retries
func GetRetryFromContext(r *http.Request) int {
	if retry, ok := r.Context().Value(retry).(int); ok {
		return retry
	}
	return 0
}

// GetAttemptsFromContext : takes the number of attempts and deduces the number of attempts already made connecting
func GetAttemptsFromContext(r *http.Request) int {
	if retry, ok := r.Context().Value(attempts).(int); ok {
		return retry
	}
	return 0
}

func doHealthCheck() {
	t := time.NewTicker(time.Second * 20)
	for {
		select {
		case <-t.C:
			log.Println("Starting health check...")
			s.HealthCheck()
			log.Println("Health check completed")
		}
	}
}

func (c *config) getConf(configFile string) {
	yamlFile, err := ioutil.ReadFile(configFile)
	if err != nil {
		log.Fatal("Could not read config file: ", err)
	}

	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		log.Fatal("Could not parse yaml: ", err)
	}
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func main() {
	var configFile string
	var c config

	flag.StringVar(&configFile, "config", "config.yml", "The name of the config file to use")
	flag.Parse()

	c.getConf(configFile)
	serverList := c.Backends
	port := c.Port
	var seenToks []string

	for _, tok := range serverList {
		if stringInSlice(tok, seenToks) {
			log.Fatal("This URL is repeated in config: ", tok)
		} else {
			seenToks = append(seenToks, tok)
		}

		serverURL, err := url.Parse(tok)
		if err != nil {
			log.Fatal("URL is malformed!")
			return
		}

		proxy := httputil.NewSingleHostReverseProxy(serverURL)
		proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, e error) {
			log.Printf("[%s] %s\n", serverURL, e.Error())
			retries := GetRetryFromContext(r)

			if retries < 3 {
				select {
				case <-time.After(10 * time.Millisecond):
					ctx := context.WithValue(r.Context(), retry, retries+1)
					proxy.ServeHTTP(w, r.WithContext(ctx))
				}
				return
			}

			s.MarkBackendStatus(serverURL, false)
			attempts := GetAttemptsFromContext(r)
			log.Printf("%s(%s) Attempting retry %d\n", r.RemoteAddr, r.URL.Path, attempts)
			ctx := context.WithValue(r.Context(), attempts, attempts+1)
			lb(w, r.WithContext(ctx))
		}
		backend := &Backend{
			URL:          serverURL,
			Alive:        true,
			ReverseProxy: proxy,
		}

		s.backends = append(s.backends, backend)
	}

	go doHealthCheck()

	server := http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: http.HandlerFunc(lb),
	}

	err := server.ListenAndServe()
	if err != nil {
		log.Fatal("An error occurred: ", err)
	}
}
