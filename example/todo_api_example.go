package example

import (
	"fmt"

	"github.com/davidp-ro/homebrew-http-server/server"
)

// Wrapper function to start the example TODOs REST API
//
// To run: `go run todo_api_example.go` from root!
func StartExampleAPI() {
	s := server.HTTPServer{Cors: server.GetDefaultCorsOptions()}
	s.On(server.Get("/", func(data server.HandlerData) []byte {
		return s.RespondWith(200, "OK!")
	}))
	s.On(server.Get("/todos", GetTodos))
	s.On(server.Get("/todos/$id", GetTodo))
	s.On(server.Post("/todos", CreateTodo))
	s.On(server.Delete("/todos/$id", DeleteTodo))
	s.Start(":8080", func() {
		fmt.Printf("TODOs Server started on localhost:8080\n")
	})
}
