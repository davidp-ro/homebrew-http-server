package http

import "fmt"

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
	return fmt.Sprintf("%s: %s", r.Method, r.URL)
}

type RequestBody struct {
	Text string
	JSON map[string]interface{}
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

	return "RequestBody{ <EMPTY> }"
}
