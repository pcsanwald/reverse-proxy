package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/nyaruka/phonenumbers"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/mail"
	"net/url"
	"os"
	"strings"
	"unicode/utf8"
)

// Define a struct to handle our configuration file format

type Configuration struct {
	Server string `json:"server"`
	Rules  Deny   `json:"deny"`
}
type Deny struct {
	Headers   []string `json:"headers"`
	URLParams []string `json:"url-params"`
}

func main() {
	configFileName := "config.json"
	if len(os.Args) > 1 {
		configFileName = os.Args[1]
	}
	body, err := os.ReadFile(configFileName)
	if err != nil {
		log.Fatalf("Unable to load configuration file, %v. Please specify as an argument to the program.", err)
	}
	reverseProxyConfig := parseConfigFile(body)
	fmt.Println("started reverse proxy...")

	proxy, err := NewProxy(reverseProxyConfig)
	if err != nil {
		log.Fatal(err)
	}
	http.HandleFunc("/", ProxyRequestHandler(proxy))
	log.Fatal(http.ListenAndServe(":9090", nil))
}

func parseConfigFile(configBytes []byte) *Configuration {
	reverseProxyConfig := Configuration{}
	err := json.Unmarshal(configBytes, &reverseProxyConfig)
	if err != nil {
		log.Fatal(err)
	}
	return &reverseProxyConfig
}

// Implement RoundTrip function so that we can log the response to the request
// from our reverse proxy, as well as potentially block requests

type loggingTransport struct {
	config *Configuration
}

func shouldBlockRequest(request *http.Request, config *Configuration) bool {
	for headerIndex := 0; headerIndex < len(config.Rules.Headers); headerIndex++ {
		headerToBlock := config.Rules.Headers[headerIndex]
		if request.Header.Get(headerToBlock) != "" {
			return true
		}
	}

	for queryParamIndex := 0; queryParamIndex < len(config.Rules.URLParams); queryParamIndex++ {
		queryParamToBlock := config.Rules.URLParams[queryParamIndex]
		if request.URL.Query().Has(queryParamToBlock) {
			return true
		}
	}
	return false
}

func looksLikeEmail(parameterValue string) bool {
	_, err := mail.ParseAddress(parameterValue)
	return err == nil
}

func looksLikePhone(parameterValue string) bool {
	potentialPhoneNumber, err := phonenumbers.Parse(parameterValue, "US")
	if err != nil {
		return false
	}
	return phonenumbers.IsValidNumber(potentialPhoneNumber)
}

func maskValue(value string) string {
	return strings.Repeat("X", utf8.RuneCountInString(value))
}
func maskQueryParameters(requestParams url.Values) url.Values {
	maskedValues := url.Values{}
	for key, values := range requestParams {
		for _, value := range values {
			if looksLikeEmail(value) || looksLikePhone(value) {
				maskedValues.Add(key, maskValue(value))
			} else {
				maskedValues.Add(key, value)
			}
		}
	}
	return maskedValues
}

func (t *loggingTransport) RoundTrip(request *http.Request) (*http.Response, error) {
	var response *http.Response
	var err error
	if shouldBlockRequest(request, t.config) {
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
		maskedValues := maskQueryParameters(request.URL.Query())
		request.URL.RawQuery = maskedValues.Encode()
		log.Printf("Query String after masking values: %v", request.URL.RawQuery)
		response, err = http.DefaultTransport.RoundTrip(request)
	}

	// 3. The proxy should log all incoming requests, including headers and body, and response
	// headers and body.
	log.Println(ResponseAsLoggableString(response))
	return response, err
}

// Creates a reverse proxy mapped to a targetHost

func NewProxy(config *Configuration) (*httputil.ReverseProxy, error) {
	url, err := url.Parse(config.Server)
	if err != nil {
		return nil, err
	}
	proxy := httputil.NewSingleHostReverseProxy(url)
	proxy.Transport = &loggingTransport{config: config}
	return proxy, nil
}

// Creates a request as formatted string, this function is meant to separate formatting
// concerns separate from logging i/o

func RequestAsLoggableString(request *http.Request) string {
	logLine, _ := httputil.DumpRequest(request, true)
	return string(logLine)
}

// Creates a response as a formatted string, this function is meant to separate formatting
// concerns separate from logging i/o

func ResponseAsLoggableString(response *http.Response) string {
	logLine, _ := httputil.DumpResponse(response, true)
	return string(logLine)
}

// Use a custom handler so that we can log the request submitted to our reverse proxy

func ProxyRequestHandler(proxy *httputil.ReverseProxy) func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		// 3. The proxy should log all incoming requests, including headers and body, and response
		// headers and body.
		// NOTE: this logging happens prior to masking PII: I want to include logging here because it's
		// part of the requirements, but in practice, we would probably not want to log PII prior to masking.
		log.Println(RequestAsLoggableString(request))
		proxy.ServeHTTP(writer, request)

	}
}
