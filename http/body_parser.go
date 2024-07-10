package http

import (
	"encoding/json"
	"fmt"
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

	switch s.ToLower(contentType) {
	case "plain/text":
		return parseText(reqString)
	case "application/json":
		return parseJSON(reqString)
	case "application/x-www-form-urlencoded":
		return parseForm(req)
	default:
		return RequestBody{}, fmt.Errorf("unsupported content type")
	}
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

func parseForm(req []byte) (RequestBody, error) {
	panic("unimplemented")
}
