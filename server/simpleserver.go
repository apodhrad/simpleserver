package server

import (
	"context"
	"fmt"
	"io/fs"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

const DEFAULT_PORT int = 8883

const CTX_NAME_KEY string = "ctx_name"

const CONTENT_TYPE_TEXT string = "text/plain"
const CONTENT_TYPE_JSON string = "application/json"
const CONTENT_TYPE_FILE string = "application/octet-stream"

// Very simple http server which makes your file directory available via http.
//
// Default port is 8883.
type SimpleServer struct {
	Port      int
	Dir       string
	server    *http.Server
	ctx       context.Context
	cancelCtx context.CancelFunc
}

// Starts the http server and return its address.
func (s *SimpleServer) Start() (string, error) {
	mux := http.NewServeMux()

	fileInfo, err := os.Stat(s.Dir)
	if err != nil {
		return "", err
	}

	if !fileInfo.IsDir() {
		return "", fmt.Errorf("server dir '%s' is not a dir", s.Dir)
	}

	err = filepath.Walk(s.Dir, func(path string, info fs.FileInfo, err error) error {
		if !info.IsDir() {
			mux.HandleFunc(simpleHandler(s.Dir, path))
		}
		return err
	})
	if err != nil {
		return "", err
	}

	s.ctx, s.cancelCtx = context.WithCancel(context.Background())
	s.server = &http.Server{
		Addr:    getAddr(s.Port),
		Handler: mux,
		BaseContext: func(l net.Listener) context.Context {
			return s.ctx
		},
	}
	go func() {
		s.server.ListenAndServe()
		s.cancelCtx()
	}()

	time.Sleep(1 * time.Second)

	return s.server.Addr, err
}

// Stops the http server
func (simpleServer *SimpleServer) Stop() error {
	err := simpleServer.server.Shutdown(simpleServer.ctx)
	if err != nil {
		return err
	}
	<-simpleServer.ctx.Done()
	return nil
}

func getAddr(port int) string {
	if port == 0 {
		port = DEFAULT_PORT
	}
	return fmt.Sprintf("localhost:%d", port)
}

func simpleHandler(dir string, path string) (string, func(w http.ResponseWriter, r *http.Request)) {
	simpleHandler := func(w http.ResponseWriter, r *http.Request) {
		fileBytes, err := os.ReadFile(path)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.Header().Set("Content-Type", CONTENT_TYPE_TEXT)
			w.WriteHeader(http.StatusOK)
			w.Write(fileBytes)
		}
	}
	pattern := path[len(dir):]
	return pattern, simpleHandler
}
