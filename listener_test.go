package main

import (
	"net/http"
	"regexp"
	"testing"
)

func TestRequestAsLoggableString(t *testing.T) {
	// TODO: this can be much more exhaustive
	request := http.Request{
		Method: "GET",
	}
	expectedOutputRegex := "Method: GET"
	validLoggableString := regexp.MustCompile(expectedOutputRegex)
	logString := RequestAsLoggableString(&request)
	if !validLoggableString.MatchString(logString) {
		t.Fatalf(`Output of RequestAsLoggableString, %s, didn't match %s`, logString, expectedOutputRegex)
	}
}
