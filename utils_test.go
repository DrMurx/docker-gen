package main

import (
	"flag"
	"os"
	"testing"
)

func TestDefaultEndpoint(t *testing.T) {
	endpoint, err := getEndpoint()
	if err != nil {
		t.Fatalf("%s", err)
	}
	if endpoint != "unix:///var/run/docker.sock" {
		t.Fatalf("Expected unix:///var/run/docker.sock, got %s", endpoint)
	}
}

func TestDockerHostEndpoint(t *testing.T) {
	err := os.Setenv("DOCKER_HOST", "tcp://127.0.0.1:4243")
	if err != nil {
		t.Fatalf("Unable to set DOCKER_HOST: %s", err)
	}

	endpoint, err := getEndpoint()
	if err != nil {
		t.Fatal("%s", err)
	}

	if endpoint != "tcp://127.0.0.1:4243" {
		t.Fatalf("Expected tcp://127.0.0.1:4243, got %s", endpoint)
	}
}

func TestDockerFlagEndpoint(t *testing.T) {

	initFlags()
	err := os.Setenv("DOCKER_HOST", "tcp://127.0.0.1:4243")
	if err != nil {
		t.Fatalf("Unable to set DOCKER_HOST: %s", err)
	}

	// flag value should override DOCKER_HOST and default value
	err = flag.Set("endpoint", "tcp://127.0.0.1:5555")
	if err != nil {
		t.Fatalf("Unable to set endpoint flag: %s", err)
	}

	endpoint, err := getEndpoint()
	if err != nil {
		t.Fatal("%s", err)
	}
	if endpoint != "tcp://127.0.0.1:5555" {
		t.Fatalf("Expected tcp://127.0.0.1:5555, got %s", endpoint)
	}
}

func TestUnixBadFormat(t *testing.T) {
	endpoint = "unix:/var/run/docker.sock"
	_, err := getEndpoint()
	if err == nil {
		t.Fatal("endpoint should have failed")
	}
}

func TestSplitKeyValueSlice(t *testing.T) {
	tests := []struct {
		input    []string
		expected string
	}{
		{[]string{"K"}, ""},
		{[]string{"K="}, ""},
		{[]string{"K=V3"}, "V3"},
		{[]string{"K=V4=V5"}, "V4=V5"},
	}

	for _, i := range tests {
		v := splitKeyValueSlice(i.input)
		if v["K"] != i.expected {
			t.Fatalf("expected K='%s'. got '%s'", i.expected, v["K"])
		}

	}
}
