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

const CTX_NAME_KEY string = "ctx_name"

const CONTENT_TYPE_TEXT string = "text/plain"
const CONTENT_TYPE_JSON string = "application/json"
const CONTENT_TYPE_FILE string = "application/octet-stream"

func SimpleHandler(dir string, path string) (string, func(w http.ResponseWriter, r *http.Request)) {
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

type SimpleServer struct {
	server    *http.Server
	ctx       context.Context
	cancelCtx context.CancelFunc
}

func (simpleServer *SimpleServer) Start(dir string) error {
	mux := http.NewServeMux()

	err := filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			mux.HandleFunc(SimpleHandler(dir, path))
		}
		return nil
	})
	if err != nil {
		return err
	}

	simpleServer.ctx, simpleServer.cancelCtx = context.WithCancel(context.Background())
	simpleServer.server = &http.Server{
		Addr:    ":8333",
		Handler: mux,
		BaseContext: func(l net.Listener) context.Context {
			// ctx = context.WithValue(ctx, CTX_NAME_KEY, ctxName)
			return simpleServer.ctx
		},
	}
	go func() {
		fmt.Printf("Start server\n")
		err := simpleServer.server.ListenAndServe()
		if errors.Is(err, http.ErrServerClosed) {
			fmt.Printf("Server closed\n")
		} else if err != nil {
			fmt.Printf("error listening for server one: %s\n", err)
		}
		simpleServer.cancelCtx()
	}()

	time.Sleep(1 * time.Second)

	return err
}

func (simpleServer *SimpleServer) Stop() {
	simpleServer.server.Shutdown(simpleServer.ctx)
	<-simpleServer.ctx.Done()
}
