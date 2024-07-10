package http

import (
	"fmt"
	"slices"
	s "strings"
)

type requestLine struct {
	Method  HTTPMethod
	URL     string
	Version string
}

func parseRequestLine(firstLine string) (requestLine, error) {
	parts := s.Split(firstLine, " ")
	if len(parts) != 3 {
		return requestLine{}, fmt.Errorf("invalid request line")
	}

	return requestLine{
		Method:  HTTPMethod(parts[0]),
		URL:     parts[1],
		Version: parts[2],
	}, nil
}

func extractPathAndQuery(url string) (string, map[string]string) {
	parts := s.Split(url, "?")
	if len(parts) == 1 {
		// No query params
		return parts[0], nil
	}

	queryParams := s.Split(parts[1], "&")
	params := make(map[string]string)

	for _, param := range queryParams {
		paramKV := s.Split(param, "=")
		if len(paramKV) != 2 {
			// Invalid
			continue
		}
		params[paramKV[0]] = paramKV[1]
	}

	return parts[0], params
}

func parseHeaders(headerLines []string) map[string]string {
	headers := make(map[string]string)
	for _, line := range headerLines {
		parts := s.Split(line, ":")
		if len(parts) < 2 {
			// invalid header
			continue
		}
		headers[parts[0]] = s.Trim(s.Join(parts[1:], ""), " ")
	}
	return headers
}

func ParseRawRequest(raw []byte) (Request, error) {
	req := string(raw[:])
	if req == "" {
		return Request{}, fmt.Errorf("empty request")
	}

	lines := s.Split(req, "\r\n")

	requestLine, err := parseRequestLine(lines[0])
	if err != nil {
		return Request{}, err
	}

	headerFields := parseHeaders(lines[1:])

	if !slices.Contains(SUPPORTED_METHODS, requestLine.Method) {
		return Request{}, fmt.Errorf("unsupported method: %s", requestLine.Method)
	}

	path, queryParams := extractPathAndQuery(requestLine.URL)

	body, err := ParseRequestIntoBodyData(headerFields["Content-Type"], req, raw)
	if err != nil {
		return Request{}, err
	}

	return Request{
		Method:      requestLine.Method,
		URL:         requestLine.URL,
		Path:        path,
		QueryParams: queryParams,
		Headers:     headerFields,
		Body:        body,
		RawRequest:  raw,
	}, nil
}
