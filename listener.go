package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
)

func main() {
	body, err := os.ReadFile("proxy.config")
	fmt.Println(string(body))

	// If config file isn't present, log but don't fail, the proxy can still
	// provide some functionality.
	if err != nil {
		fmt.Printf("unable to read configuration file: %v", err)
	}
	fmt.Println("started reverse proxy...")

	//TODO: remove hardcoded host, specify in config file
	proxy, err := NewProxy("http://localhost:8080")
	if err != nil {
		log.Fatal(err)
	}
	http.HandleFunc("/", ProxyRequestHandler(proxy))
	log.Fatal(http.ListenAndServe(":9090", nil))
}

// Implement RoundTrip function so that we can log the response to the request
// from our reverse proxy.

type loggingTransport struct {
}

func (t *loggingTransport) RoundTrip(request *http.Request) (*http.Response, error) {
	var response *http.Response
	var err error
	log.Printf("beginning %v", request)
	if strings.Contains(request.URL.Path, "/foo") {
		log.Println("Round trip: CONTAINS FOO RETURNING 403")
		response = &http.Response{
			Status:        "403 Forbidden",
			StatusCode:    403,
			Proto:         "HTTP/1.1",
			ProtoMajor:    1,
			ProtoMinor:    1,
			Body:          io.NopCloser(bytes.NewBufferString("")),
			ContentLength: 0,
			Request:       request,
			Header:        make(http.Header, 0),
		}
	} else {
		response, err = http.DefaultTransport.RoundTrip(request)
	}
	
	// 3. The proxy should log all incoming requests, including headers and body, and response
	// headers and body.
	body, err := httputil.DumpResponse(response, true)
	log.Printf("Response: %v", string(body))
	if err != nil {
		return nil, err
	}
	return response, err
}

// Creates a reverse proxy mapped to a targetHost

func NewProxy(targetHost string) (*httputil.ReverseProxy, error) {
	url, err := url.Parse(targetHost)
	if err != nil {
		return nil, err
	}
	proxy := httputil.NewSingleHostReverseProxy(url)
	proxy.Transport = &loggingTransport{}
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
