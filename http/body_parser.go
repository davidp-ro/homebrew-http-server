package http

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	s "strings"
)

func ParseRequestIntoBodyData(
	contentType string,
	reqString string,
	req []byte,
) (RequestBody, error) {
	if contentType == "" {
		return RequestBody{}, nil
	}
	contentType = s.Trim(contentType, " ")

	if s.HasPrefix(contentType, "multipart") {
		return parseMultipart(req, contentType)
	}

	switch s.ToLower(contentType) {
	case "plain/text":
		return parseText(reqString)
	case "application/json":
		return parseJSON(reqString)
	case "application/x-www-form-urlencoded":
		return parseForm(reqString)
	}

	return RequestBody{}, fmt.Errorf("unsupported content type")
}

func parseText(reqString string) (RequestBody, error) {
	sections := s.Split(reqString, "\r\n\r\n")
	if len(sections) < 2 {
		return RequestBody{}, nil
	}
	return RequestBody{Text: s.Join(sections[1:], "\n\n")}, nil
}

func parseJSON(reqString string) (RequestBody, error) {
	sections := s.Split(reqString, "\r\n\r\n")
	if len(sections) < 2 {
		return RequestBody{}, nil
	}

	jsonString := s.Join(sections[1:], "\n\n")
	jsonString = s.ReplaceAll(jsonString, "\x00", "")
	jsonData := make(map[string]interface{})
	err := json.Unmarshal([]byte(jsonString), &jsonData)

	return RequestBody{JSON: jsonData}, err
}

func parseForm(reqString string) (RequestBody, error) {
	sections := s.Split(reqString, "\r\n\r\n")
	if len(sections) < 2 {
		return RequestBody{}, nil
	}
	return RequestBody{FormItems: ExtractQueryParamsFrom(sections[1])}, nil
}

func parseMultipart(req []byte, contentType string) (RequestBody, error) {
	parts := make([]MultipartBody, 1)

	_, params, err := mime.ParseMediaType(contentType)
	if err != nil {
		return RequestBody{}, fmt.Errorf("parseMultipart: %v", err)
	}

	boundary, ok := params["boundary"]
	if !ok {
		return RequestBody{}, fmt.Errorf("parseMultipart: boundary not in header")
	}

	reader := multipart.NewReader(bytes.NewReader(req), boundary)
	for {
		part, err := reader.NextPart()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return RequestBody{Multipart: parts}, nil
			}
			return RequestBody{}, fmt.Errorf("parseMultipart: %v", err)
		}

		// Read part data
		var partData bytes.Buffer
		_, err = io.Copy(&partData, part)
		if err != nil {
			return RequestBody{}, fmt.Errorf("parseMultipart: %v", err)
		}

		// Process each part based on its Content-Disposition
		disposition := part.Header.Get("Content-Disposition")
		if disposition != "" {
			_, params, err := mime.ParseMediaType(disposition)
			if err != nil {
				return RequestBody{}, fmt.Errorf("parseMultipart: %v", err)
			}
			name := params["name"]
			filename := params["filename"]

			if filename != "" {
				// Handle a file
				parts = append(parts, MultipartBody{
					Filename: filename,
					Name:     name,
					Headers:  part.Header,
					FileData: partData,
				})
			} else {
				// Handle form field
				parts = append(parts, MultipartBody{
					Name:     name,
					Headers:  part.Header,
					FormData: partData.String(),
				})
			}
		}
	}
}
