package server

import (
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

const STATUS_200_OK string = "200 OK"
const STATUS_404_NOT_FOUND string = "404 Not Found"

func assertRequest(t *testing.T, addr string, path string, expected string) {
	url := "http://" + addr + path
	resp, err := http.Get(url)
	assert.Nil(t, err)
	if resp != nil {
		defer resp.Body.Close()
	}

	actual, err := io.ReadAll(resp.Body)
	assert.Nil(t, err)
	assert.Equal(t, expected, string(actual))
}

func assertStatus(t *testing.T, addr string, path string, expectedStatus string) {
	url := "http://" + addr + path
	resp, err := http.Get(url)
	if resp != nil {
		defer resp.Body.Close()
		assert.Equal(t, resp.Status, expectedStatus)
	} else {
		assert.Fail(t, err.Error())
	}
}

func assertRequestError(t *testing.T, addr string, path string) {
	url := "http://" + addr + path
	resp, err := http.Get(url)
	assert.NotNil(t, err)
	if resp != nil {
		defer resp.Body.Close()
	}
}
func TestSimpleServer_00_empty(t *testing.T) {
	server := SimpleServer{Dir: "test-resources/00_empty"}
	addr, err := server.Start()

	assert.Nil(t, err)
	assertStatus(t, addr, "", STATUS_404_NOT_FOUND)
	assertStatus(t, addr, "/hello.txt", STATUS_404_NOT_FOUND)
	assertStatus(t, addr, "/hello/helloworld.txt", STATUS_404_NOT_FOUND)

	server.Stop()
	assertRequestError(t, addr, "")
}

func TestSimpleServer_01_helloworld(t *testing.T) {
	server := SimpleServer{Dir: "test-resources/01_helloworld"}
	addr, err := server.Start()

	assert.Nil(t, err)
	assertStatus(t, addr, "", STATUS_404_NOT_FOUND)
	assertStatus(t, addr, "/hello.txt", STATUS_200_OK)
	assertRequest(t, addr, "/hello.txt", "Hello SimpleServer!")
	assertStatus(t, addr, "/hello/helloworld.txt", STATUS_200_OK)
	assertRequest(t, addr, "/hello/helloworld.txt", "Hello World!")

	server.Stop()
	assertRequestError(t, addr, "")
}

func TestSimpleServer_02_multi(t *testing.T) {
	server1 := SimpleServer{Port: 8881, Dir: "test-resources/02_multi/01_sub"}
	addr1, err1 := server1.Start()

	server2 := SimpleServer{Port: 8882, Dir: "test-resources/02_multi/02_sub"}
	addr2, err2 := server2.Start()

	assert.Nil(t, err1)
	assert.Nil(t, err2)
	assertRequest(t, addr1, "/hello.txt", "Hello from 01_sub")
	assertRequest(t, addr2, "/hello.txt", "Hello from 02_sub")

	server1.Stop()
	assertRequestError(t, addr1, "")
	// server2 should be still running
	assertRequest(t, addr2, "/hello.txt", "Hello from 02_sub")
	server2.Stop()
	assertRequestError(t, addr2, "")
}
