package server

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/davidp-ro/homebrew-http-server/http"
	"github.com/davidp-ro/homebrew-http-server/tcp"
	"github.com/davidp-ro/homebrew-http-server/utils"
)

type HandlerFunc func(HandlerData) []byte

type HandlerData struct {
	Request    http.Request
	PathParams map[string]string
}

type Handler struct {
	OnMethod http.HTTPMethod
	Path     string
	Handler  HandlerFunc
}

type CorsOptions struct {
	AllowOrigin  string
	AllowMethods string
	AllowHeaders string
}

type HTTPServer struct {
	Handlers []Handler
	Cors     CorsOptions
	Debug    bool
}

func match(s HTTPServer, req http.Request) []byte {
	if s.Debug {
		fmt.Printf("[D] HTTPServer::match - Request: %s\n", req)
	}

	methodHandlers := utils.Filter(s.Handlers, func(h Handler) bool {
		return h.OnMethod == req.Method
	})
	if len(methodHandlers) == 0 {
		if req.Method != http.OPTIONS {
			// For all methods except OPTIONS, return 405
			return s.RespondWith(405, "Method Not Allowed")
		} else {
			// Otherwise, handle CORS preflight request (if CORS is enabled)
			if s.Cors != (CorsOptions{}) {
				return s.RespondWith(204, "")
			}
		}
	}

	// Look for exact path match:
	pathHandlers := utils.Filter(methodHandlers, func(h Handler) bool {
		return h.Path == req.Path
	})
	if len(pathHandlers) > 1 {
		fmt.Printf("HTTPServer::match - Multiple handlers for the same path: %s\n", req.Path)
		return s.RespondWith(500, "Internal Server Error")
	}
	if len(pathHandlers) == 0 && req.Method == http.OPTIONS {
		// Handle CORS preflight request (if CORS is enabled) implicitly
		if s.Cors != (CorsOptions{}) {
			return s.RespondWith(204, "")
		}
	}

	pathParams := make(map[string]string)
	if len(pathHandlers) == 0 {
		if s.Debug {
			fmt.Printf("[D] HTTPServer::match - No exact path match for: %s\n", req.Path)
		}

		// Try to find a path with path params
		params, handler := matchToParametrizedPath(s, methodHandlers, req)
		if params == nil {
			return s.RespondWith(404, "Not Found")
		}

		pathParams = params
		pathHandlers = []Handler{handler}
	}

	return pathHandlers[0].Handler(HandlerData{
		Request: req, PathParams: pathParams,
	})
}

func (s *HTTPServer) On(handler Handler) HTTPServer {
	s.Handlers = append(s.Handlers, handler)
	if s.Debug {
		fmt.Printf("[D] HTTPServer::On - Added handler: %v\n", handler)
	}
	return *s
}

func newHandler(m http.HTTPMethod, p string, d HandlerFunc) Handler {
	return Handler{OnMethod: m, Path: p, Handler: d}
}

func Get(path string, handler HandlerFunc) Handler {
	return newHandler(http.GET, path, handler)
}

func Post(path string, handler HandlerFunc) Handler {
	return newHandler(http.POST, path, handler)
}

func Put(path string, handler HandlerFunc) Handler {
	return newHandler(http.PUT, path, handler)
}

func Delete(path string, handler HandlerFunc) Handler {
	return newHandler(http.DELETE, path, handler)
}

func Options(path string, handler HandlerFunc) Handler {
	return newHandler(http.OPTIONS, path, handler)
}

func (s HTTPServer) RespondWith(statusCode int, body string, customHeaders ...map[string]string) []byte {
	statusTexts := map[int]string{
		200: "OK",
		201: "Created",
		204: "No Content",
		400: "Bad Request",
		401: "Unauthorized",
		403: "Forbidden",
		404: "Not Found",
		405: "Method Not Allowed",
		500: "Internal Server Error",
	}

	statusText := statusTexts[statusCode]
	if statusText == "" {
		panic("HTTPServer::RespondWith: Unsupported status code")
	}

	headers := GetDefaultHeaders()

	// Add (or overwrite with) custom headers
	for _, customHeader := range customHeaders {
		for k, v := range customHeader {
			headers[strings.ToLower(k)] = v
		}
	}

	// Add CORS headers
	for k, v := range GetCorsHeaders(s) {
		headers[strings.ToLower(k)] = v
	}

	return []byte(
		"HTTP/1.1 " + strconv.Itoa(statusCode) + " " + statusText + "\r\n" +
			BuildHeadersString(headers) + "\r\n" +
			body,
	)
}

func (s HTTPServer) Start(address string, onStart func()) {
	tcp.Server{
		Network: tcp.NetworkTCP,
		Address: address,
		OnConnection: func(data []byte, err error) []byte {
			if err != nil {
				log.Printf("HTTPServer::OnConnection - Fail: %v\n", err)
				return s.RespondWith(500, "Internal Server Error")
			}

			parsed, err := http.ParseRawRequest(data)
			if err != nil {
				log.Printf("HTTPServer::OnConnection - Parsing Error: %v\n", err)
				return s.RespondWith(400, "Bad Request")
			}

			return match(s, parsed)
		},
	}.Start(onStart)
}
