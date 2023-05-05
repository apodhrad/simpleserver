package server

import (
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

const STATUS_200_OK string = "200 OK"
const STATUS_404_NOT_FOUND string = "404 Not Found"

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

func assertStatus(t *testing.T, url string, expectedStatus string) {
	resp, err := http.Get(url)
	if resp != nil {
		defer resp.Body.Close()
		assert.Equal(t, resp.Status, expectedStatus)
	} else {
		assert.Fail(t, err.Error())
	}
}

func assertRequestError(t *testing.T, url string) {
	resp, err := http.Get(url)
	assert.NotNil(t, err)
	if resp != nil {
		defer resp.Body.Close()
	}
}
func TestSimpleServer_00_empty(t *testing.T) {
	server := SimpleServer{}
	server.Start("test-resources/00_empty")

	assertStatus(t, "http://localhost:8333", STATUS_404_NOT_FOUND)
	assertStatus(t, "http://localhost:8333/hello.txt", STATUS_404_NOT_FOUND)
	assertStatus(t, "http://localhost:8333/hello/helloworld.txt", STATUS_404_NOT_FOUND)

	server.Stop()
	assertRequestError(t, "http://localhost:8333")
}

func TestSimpleServer_01_helloworld(t *testing.T) {
	server := SimpleServer{}
	server.Start("test-resources/01_helloworld")

	assertStatus(t, "http://localhost:8333", STATUS_404_NOT_FOUND)
	assertStatus(t, "http://localhost:8333/hello.txt", STATUS_200_OK)
	assertRequest(t, "http://localhost:8333/hello.txt", "Hello SimpleServer!")
	assertStatus(t, "http://localhost:8333/hello/helloworld.txt", STATUS_200_OK)
	assertRequest(t, "http://localhost:8333/hello/helloworld.txt", "Hello World!")

	server.Stop()
	assertRequestError(t, "http://localhost:8333")
}
