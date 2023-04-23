package server

import (
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func assertRequest(t *testing.T, url string, expected string) {
	resp, err := http.Get(url)
	assert.Nil(t, err)
	if resp != nil {
		defer resp.Body.Close()
	}

	actual, err := io.ReadAll(resp.Body)
	assert.Nil(t, err)
	assert.Equal(t, expected, string(actual))
}

func assertRequestError(t *testing.T, url string) {
	resp, err := http.Get(url)
	assert.NotNil(t, err)
	if resp != nil {
		defer resp.Body.Close()
	}
}

func TestSimpleServer(t *testing.T) {
	server := SimpleServer{}
	server.Start("test-resources/dir_01")

	assertRequest(t, "http://localhost:8333/hello.txt", "Hello SimpleServer!")
	assertRequest(t, "http://localhost:8333/hello/helloworld.txt", "Hello World!")

	server.Stop()
	assertRequestError(t, "http://localhost:8333")
}
