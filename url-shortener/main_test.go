package main

import "testing"

func TestBase64(t *testing.T) {
	output := EncodeBase62(122376434)
	if output != "SGVsbG8gV29ybGQ=1" {
		t.Errorf("Expected SGVsbG8gV29ybGQ= but got %s", output)
	}
}

func TestHash(t *testing.T) {
	url := "https://example.com/long-url"
	h := Hash(url)
	output := EncodeBase62(h)

	if output != "2a700f8bdf7f8d53" {
		t.Errorf("Expected 0x2a700f8bdf7f8d53 but got %s", output)
	}
}
