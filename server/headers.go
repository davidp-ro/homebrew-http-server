package server

// Builds the headers string (as expected in a HTTP Response) from a map of
// headers. One per line, with the format "key: value\n".
func BuildHeadersString(headers map[string]string) string {
	str := ""
	for k, v := range headers {
		str += k + ": " + v + "\r\n"
	}
	return str
}

// Returns a map of CORS headers from a HTTPServer struct. Assumes that CORS
// is setup on the server.
func GetCorsHeaders(s HTTPServer) map[string]string {
	return map[string]string{
		"access-control-allow-origin":  s.Cors.AllowOrigin,
		"access-control-allow-methods": s.Cors.AllowMethods,
		"access-control-allow-headers": s.Cors.AllowHeaders,
	}
}

// Get some default (lax) CORS options.
func GetDefaultCorsOptions() CorsOptions {
	return CorsOptions{
		AllowOrigin:  "*",
		AllowMethods: "GET, POST, PUT, DELETE, OPTIONS",
		AllowHeaders: "content-type",
	}
}

// Get the default headers for a HTTP response.
func GetDefaultHeaders() map[string]string {
	return map[string]string{
		"content-type": "text/plain",
		"server":       "Homebrew HTTP Server",
		"connection":   "close",
	}
}
