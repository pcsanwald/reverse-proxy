package main

import (
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"regexp"
	"testing"
)

func TestRequestAsLoggableString(t *testing.T) {
	// TODO: this can be much more exhaustive
	testURL, _ := url.Parse("/foo")
	request := http.Request{
		Method:           "GET",
		URL:              testURL,
		Proto:            "HTTP/1.1",
		ProtoMajor:       1,
		ProtoMinor:       1,
		Header:           make(http.Header, 0),
		Body:             nil,
		GetBody:          nil,
		ContentLength:    0,
		TransferEncoding: nil,
		Close:            false,
		Host:             "localhost:9090",
		Form:             nil,
		PostForm:         nil,
		MultipartForm:    nil,
		Trailer:          nil,
		RemoteAddr:       "",
		RequestURI:       "",
		TLS:              nil,
		Response:         nil,
	}
	expectedOutputRegex := "GET /foo HTTP/1.1"
	validLoggableString := regexp.MustCompile(expectedOutputRegex)
	logString := RequestAsLoggableString(&request)
	if !validLoggableString.MatchString(logString) {
		t.Fatalf(`Output of RequestAsLoggableString, %s, didn't match %s`, logString, expectedOutputRegex)
	}
}

func TestReverseProxyBlockingByHeader(t *testing.T) {
	backendServer := httptest.NewServer(http.DefaultServeMux)
	defer backendServer.Close()
	backendURL, err := url.Parse(backendServer.URL)
	if err != nil {
		t.Fatal(err)
	}

	config := Configuration{
		Server: backendURL.String(),
		Rules: Deny{
			// Use a header that should result in request being blocked, since http.get includes
			// User-Agent by default.
			Headers:   []string{"User-Agent"},
			URLParams: []string{},
		},
	}
	reverseProxy, err := NewProxy(&config)
	if err != nil {
		log.Fatal(err)
	}

	testServer := httptest.NewServer(reverseProxy)
	defer testServer.Close()

	response, err := http.Get(testServer.URL)
	if err != nil {
		t.Fatal(err)
	}

	if response.StatusCode != 403 {
		t.Fatalf("Expecting a 403 forbidden response, got %v", response)
	}

}

// TODO: remove duplication across 2 tests?

func TestReverseProxyBlockingByParam(t *testing.T) {
	backendServer := httptest.NewServer(http.DefaultServeMux)
	defer backendServer.Close()
	backendURL, err := url.Parse(backendServer.URL)
	if err != nil {
		t.Fatal(err)
	}

	paramToBlock := "password"
	config := Configuration{
		Server: backendURL.String(),
		Rules: Deny{
			// Use a URL param for blocking
			Headers:   []string{},
			URLParams: []string{paramToBlock},
		},
	}
	reverseProxy, err := NewProxy(&config)
	if err != nil {
		log.Fatal(err)
	}

	testServer := httptest.NewServer(reverseProxy)
	defer testServer.Close()

	response, err := http.Get(testServer.URL + "?" + paramToBlock + "=asdf")
	if err != nil {
		t.Fatal(err)
	}

	if response.StatusCode != 403 {
		t.Fatalf("Expecting a 403 forbidden response, got %v", response)
	}

}
