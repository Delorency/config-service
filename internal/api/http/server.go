package http

import (
	"fmt"
	"log"
	"net/http"
)

func NewHTTPServer(addr string, port int, logger *log.Logger, schemadir string) *http.Server {
	router := NewRouter(logger, schemadir)

	server := http.Server{
		Addr:    fmt.Sprintf("%s:%d", addr, port),
		Handler: router,
	}

	return &server
}
