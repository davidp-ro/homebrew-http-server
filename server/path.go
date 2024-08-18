package server

import (
	"fmt"
	"strings"

	"github.com/davidp-ro/homebrew-http-server/http"
)

func extractPathParams(matcherPath string, reqPath string) map[string]string {
	matcherParts := strings.Split(matcherPath, "/")
	reqParts := strings.Split(reqPath, "/")

	if len(matcherParts) != len(reqParts) {
		return nil
	}

	params := make(map[string]string)
	for i, part := range matcherParts {
		if part != "" {
			params[part] = reqParts[i]
		}
	}

	return params
}

func checkPathParams(pathParams map[string]string) bool {
	if len(pathParams) == 0 {
		return false
	}

	for param := range pathParams {
		// If a matched "param" isn't dynamic (doesn't start with "$") and
		// doesn't match the path exactly, return false
		if !strings.Contains(param, "$") && param != pathParams[param] {
			return false
		}
	}

	return true
}

func matchToParametrizedPath(s HTTPServer, methodHandlers []Handler, req http.Request) (map[string]string, Handler) {
	for _, handler := range methodHandlers {
		if strings.Contains(handler.Path, "$") {
			pathParams := extractPathParams(handler.Path, req.Path)

			if s.Debug {
				fmt.Printf("[D] HTTPServer::matchToParametrizedPath - Path params: %v\n", pathParams)
			}

			if pathParams != nil {
				if checkPathParams(pathParams) {
					return pathParams, handler
				} else {
					if s.Debug {
						fmt.Printf("[D] HTTPServer::matchToParametrizedPath - Static path params don't match the path: %v\n", pathParams)
					}
				}
			}
		}
	}

	return nil, Handler{}
}
