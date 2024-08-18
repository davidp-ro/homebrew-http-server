package server

import "testing"

func TestBuildHeadersString(t *testing.T) {
	headers := map[string]string{
		"content-type":     "text/plain",
		"server":           "Homebrew HTTP Server",
		"connection":       "close",
		"some-header-here": "some-value-here",
	}
	expected := "content-type: text/plain\r\nserver: Homebrew HTTP Server\r\nconnection: close\r\nsome-header-here: some-value-here\r\n"
	if BuildHeadersString(headers) != expected {
		t.Errorf("\nExpected:\n%s=====\nGot:\n%s", expected, BuildHeadersString(headers))
	}
}
