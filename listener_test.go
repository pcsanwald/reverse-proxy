package main

import (
	"net/http"
	"regexp"
	"testing"
)

func TestRequestAsLoggableString(t *testing.T) {
	request := http.Request{
		Method:           "GET",
	}
	expectedOutputRegex := "Method: GET"
	validLoggableString := regexp.MustCompile(expectedOutputRegex)
	logString := RequestAsLoggableString(&request)
	if !validLoggableString.MatchString(logString) {
		t.Fatalf(`Output of RequestAsLoggableString, %s, didn't match %s`, logString, expectedOutputRegex)
	}
}
