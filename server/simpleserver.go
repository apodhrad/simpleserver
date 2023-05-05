package server

import (
	"context"
	"errors"
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

type SimpleServer struct {
	Port      int
	Dir       string
	server    *http.Server
	ctx       context.Context
	cancelCtx context.CancelFunc
}

func (s *SimpleServer) Start() (string, error) {
	mux := http.NewServeMux()

	err := filepath.Walk(s.Dir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			mux.HandleFunc(simpleHandler(s.Dir, path))
		}
		return nil
	})
	if err != nil {
		return "", err
	}

	s.ctx, s.cancelCtx = context.WithCancel(context.Background())
	s.server = &http.Server{
		Addr:    getAddr(s.Port),
		Handler: mux,
		BaseContext: func(l net.Listener) context.Context {
			// ctx = context.WithValue(ctx, CTX_NAME_KEY, ctxName)
			return s.ctx
		},
	}
	go func() {
		fmt.Printf("Start server\n")
		err := s.server.ListenAndServe()
		if errors.Is(err, http.ErrServerClosed) {
			fmt.Printf("Server closed\n")
		} else if err != nil {
			fmt.Printf("error listening for server one: %s\n", err)
		}
		s.cancelCtx()
	}()

	time.Sleep(1 * time.Second)

	return s.server.Addr, err
}

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
			panic(err)
		}
		w.Header().Set("Content-Type", CONTENT_TYPE_TEXT)
		w.WriteHeader(http.StatusOK)
		w.Write(fileBytes)
	}
	pattern := path[len(dir):]
	return pattern, simpleHandler
}
