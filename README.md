# homebrew-http-server

A simple, barebones HTTP/1.1 server.

My main goal for this was to better understand the protocols around the web - while new stuff like QUIC is on another level, I wanted to see how an older, simpler protocol works (especially the underlying standards, parsing, response formatting, everything really).

And I also haven't used Go in a while, figured this would be a good opportunity to get back into it!

> ✨ This is very much not meant to be a production-ready server, but rather a learning tool. I did, however, try to bake in some basic neat features:

- **Uses goroutines**: Each incoming TCP connection is handled in its own goroutine.

- **Automatic body parsing**: The server will automatically parse the body of incoming requests and make it available to the handler. Supports:

  - JSON;
  - FormData;
  - MultiPart.

- **Support for path parameters**: You can define path parameters (`/$param-part/static-part`) in your routes and they will be automatically parsed and made available to the handler.

### Packages

- `example` - an example REST API (TODOs) that uses the server.
- `http` - protocol-related code, such as parsing requests, headers, etc.
- `server` - the server implementation, includes the route/path matcher, path parameter resolving, etc.
- `tcp` - underlying TCP server implementation with the connection handling.
- `utils` - various utility functions.

### Basic Usage

```go
// Create a new server instance, with some lax CORS options:
s := server.HTTPServer{Cors: server.GetDefaultCorsOptions()}

// Example route definition:
s.On(server.Get("/todos/$id", func(data server.HandlerData) []byte {
    return s.RespondWith(200, "Hello, todo " + data.PathParams["$id"] + "!")
}))

// Start the server:
s.Start(":8080", func() {
    fmt.Println("Server started!")
})
```

> 💡 If you wanna give the API in the `example` folder a _go_: `go run todo_api_example.go`

### Limitations

- **No HTTPS support**: This server only supports HTTP.
- **No support for HTTP/2 or above**: This server only supports HTTP/1.1.
- **Many modern features are missing**: You have no compression support, you have to manually set relevant headers, performance wasn't a focus, etc.
- **My not-perfect understanding of Golang**: I'm sure there are some things that could be done better, but I really enjoyed writing some Go again!
