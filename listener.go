package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
)


// Implement RoundTrip function so that we can log the response to the request
// from our reverse proxy.

type myTransport struct {
}

func (t *myTransport) RoundTrip(request *http.Request) (*http.Response, error) {
	response, err := http.DefaultTransport.RoundTrip(request)
	body, err := httputil.DumpResponse(response, true)
	if err != nil {
		return nil, err
	}
	log.Println(string(body))
	return response, err
}

// Creates a reverse proxy mapped to a targetHost

func NewProxy(targetHost string) (*httputil.ReverseProxy, error) {
	url, err := url.Parse(targetHost)
	if err != nil {
		return nil, err
	}
	proxy := httputil.NewSingleHostReverseProxy(url)
	proxy.Transport = &myTransport{}
	return proxy, nil
}

// Creates a formatted string, this function is meant to separate formatting
// concerns separate from logging i/o

func RequestAsLoggableString(request *http.Request) string {
	logLine := fmt.Sprintf("Header: %s  Body: %s Method: %s URL: %s",
		request.Header, request.Body, request.Method, request.URL)
	return logLine
}

// Use a custom handler so that we can log the request submitted to our reverse proxy

func ProxyRequestHandler(proxy *httputil.ReverseProxy) func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		log.Println(RequestAsLoggableString(request))
		proxy.ServeHTTP(writer, request)
	}
}

func main() {
	body, err := os.ReadFile("proxy.config")
	fmt.Println(string(body))

	// If config file isn't present, log but don't fail, the proxy can still
	// provide some functionality.
	if err != nil {
		fmt.Printf("unable to read configuration file: %v", err)
	}
	fmt.Println("started reverse proxy...")

	proxy, err := NewProxy("http://localhost:8080")
	if err != nil {
		log.Fatal(err)
	}
	http.HandleFunc("/", ProxyRequestHandler(proxy))
	log.Fatal(http.ListenAndServe(":9090", nil))
}
