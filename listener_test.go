package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"regexp"
	"testing"
)

const userAgentHeaderName = "User-Agent"

func TestRequestAsLoggableString(t *testing.T) {
	// TODO: this can be much more exhaustive
	testURL, _ := url.Parse("/foo")
	request := http.Request{
		Method:        http.MethodGet,
		URL:           testURL,
		ProtoMajor:    2,
		ProtoMinor:    0,
		Header:        make(http.Header, 0),
		ContentLength: 0,
		Close:         false,
		Host:          "localhost:9090",
	}
	expectedOutputRegex := "GET /foo HTTP/2.0"
	validLoggableString := regexp.MustCompile(expectedOutputRegex)
	logString := RequestAsLoggableString(&request)
	if !validLoggableString.MatchString(logString) {
		t.Fatalf(`Output of RequestAsLoggableString, %s, didn't match %s`, logString, expectedOutputRegex)
	}
}

func TestMaskingString(t *testing.T) {
	inputString := "paul"
	expectedResult := "XXXX"
	actualResult := maskValue(inputString)
	if actualResult != expectedResult {
		t.Fatalf("for input string %v, expected %v, got %v", inputString, expectedResult, actualResult)
	}
}

func TestLooksLikeEmail(t *testing.T) {
	maybeEmail := "paul@gmail.com"
	if !looksLikeEmail(maybeEmail) {
		t.Fatalf("input string '%v' looks like an email", maybeEmail)
	}
	maybeEmail = "paul"
	if looksLikeEmail(maybeEmail) {
		t.Fatalf("input string '%v' should not look like an email", maybeEmail)
	}
}

func TestLooksLikePhone(t *testing.T) {
	maybePhone := "+33 7 69 24 58 46"
	if !looksLikePhone(maybePhone) {
		t.Fatalf("input string '%v' looks like a phone number", maybePhone)
	}
	maybePhone = "8675"
	if looksLikePhone(maybePhone) {
		t.Fatalf("input string '%v' should not look like a phone number", maybePhone)
	}
}

func TestShouldBlockRequestWithPOST(t *testing.T) {
	// We'll configure our proxy to block User-Agent, and then
	// set User-Agent in our request, to ensure that we are only
	// blocking GET requests, not POST
	request := http.Request{
		Method: http.MethodPost,
		Header: make(http.Header, 0),
	}
	request.Header.Set(userAgentHeaderName, "whatever")
	config := Configuration{
		Server: "not required for test",
		Rules: Deny{
			// Use a header that should result in request being blocked, since http.get includes
			// User-Agent by default.
			Headers:   []string{userAgentHeaderName},
			URLParams: []string{},
		},
	}
	if shouldBlockRequest(&request, &config) != false {
		t.Fatalf("We should not block a request with method %v", http.MethodPost)
	}
}

func TestShouldBlockRequestWithGET(t *testing.T) {
	// We'll configure our proxy to block User-Agent, and then
	// set User-Agent in our request, to ensure that we are only
	// blocking GET requests, not POST
	request := http.Request{
		Method: http.MethodGet,
		Header: make(http.Header, 0),
	}
	request.Header.Set(userAgentHeaderName, "whatever")
	config := Configuration{
		Server: "not required for test",
		Rules: Deny{
			// Use a header that should result in request being blocked, since http.get includes
			// User-Agent by default.
			Headers:   []string{userAgentHeaderName},
			URLParams: []string{},
		},
	}
	if shouldBlockRequest(&request, &config) != true {
		t.Fatalf("We should not block a request with method %v", http.MethodGet)
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
			Headers:   []string{userAgentHeaderName},
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

	if response.StatusCode != http.StatusForbidden {
		t.Fatalf("Expecting a %d response, got %v", http.StatusForbidden, response)
	}

}

func TestReverseProxyPassthroughRequest(t *testing.T) {
	backendServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "this call was relayed by the reverse proxy")
	}))
	defer backendServer.Close()
	backendURL, err := url.Parse(backendServer.URL)
	if err != nil {
		t.Fatal(err)
	}

	config := Configuration{
		Server: backendURL.String(),
		Rules: Deny{
			// Use a URL param for blocking
			Headers:   []string{},
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

	if response.StatusCode != http.StatusOK {
		t.Fatalf("Expecting a 200 OK response, got %v", response)
	}
}
