package main

import (
	"net/http"
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
