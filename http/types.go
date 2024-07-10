package http

import (
	"bytes"
	"fmt"
	"net/textproto"
)

type HTTPMethod string

const (
	GET     HTTPMethod = "GET"
	POST    HTTPMethod = "POST"
	PUT     HTTPMethod = "PUT"
	DELETE  HTTPMethod = "DELETE"
	OPTIONS HTTPMethod = "OPTIONS"
)

var SUPPORTED_METHODS = []HTTPMethod{GET, POST, PUT, DELETE, OPTIONS}

type Request struct {
	Method      HTTPMethod
	URL         string
	Path        string
	QueryParams map[string]string
	Headers     map[string]string
	Body        RequestBody
	RawRequest  []byte
}

func (r Request) String() string {
	return fmt.Sprintf("%s %s", r.Method, r.URL)
}

type RequestBody struct {
	Text      string
	JSON      map[string]interface{}
	FormItems map[string]string
	Multipart []MultipartBody
}

type MultipartBody struct {
	Filename string
	Name     string
	Headers  textproto.MIMEHeader
	FileData bytes.Buffer
	FormData string
}

func (rb RequestBody) String() string {
	if rb.Text != "" {
		return fmt.Sprintf("RequestBody(Text){ %s }", rb.Text)
	}

	if rb.JSON != nil {
		json := ""
		for k, v := range rb.JSON {
			json += fmt.Sprintf("\t%s: %v\n", k, v)
		}
		return fmt.Sprintf("RequestBody(JSON){\n%s}", json)
	}

	if len(rb.FormItems) > 0 {
		form := ""
		for k, v := range rb.FormItems {
			form += fmt.Sprintf("\t%s: %s\n", k, v)
		}
		return fmt.Sprintf("RequestBody(FormItems){\n%s}", form)
	}

	if len(rb.Multipart) > 0 {
		return fmt.Sprintf("RequestBody(Multipart){ Parts: %d }", len(rb.Multipart))
	}

	return "RequestBody{ <EMPTY> }"
}
