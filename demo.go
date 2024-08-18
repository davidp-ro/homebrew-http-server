package main

import (
	"fmt"

	"github.com/davidp-ro/homebrew-http-server/server"
)

func main() {
	s := server.HTTPServer{Debug: true, Cors: server.GetDefaultCorsOptions()}
	s.On(server.Get("/", func(data server.HandlerData) []byte {
		return s.RespondWith(200, "Hello!")
	}))
	s.On(server.Get("/hello", func(data server.HandlerData) []byte {
		return s.RespondWith(200, "Hello, world!")
	}))
	s.On(server.Options("/hello", func(data server.HandlerData) []byte {
		return s.RespondWith(200, "Hello, options!")
	}))
	s.On(server.Get("/$path1/$path2/test", func(data server.HandlerData) []byte {
		return s.RespondWith(200, fmt.Sprintf("Hello, path! %v", data.PathParams))
	}))
	s.On(server.Post("/$path1/test/$path2", func(data server.HandlerData) []byte {
		return s.RespondWith(200, fmt.Sprintf("Hello, path 2! %v", data.PathParams))
	}))
	s.Start(":8080", func() {
		println("Server started on :8080")
	})
}
