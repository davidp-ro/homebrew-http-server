package server

func BuildHeadersString(headers map[string]string) string {
	str := ""
	for k, v := range headers {
		str += k + ": " + v + "\n"
	}
	return str
}

func GetCorsHeaders(s HTTPServer) map[string]string {
	return map[string]string{
		"access-control-allow-origin":  s.Cors.AllowOrigin,
		"access-control-allow-methods": s.Cors.AllowMethods,
		"access-control-allow-headers": s.Cors.AllowHeaders,
	}
}

func GetDefaultCorsOptions() CorsOptions {
	return CorsOptions{
		AllowOrigin:  "*",
		AllowMethods: "GET, POST, PUT, DELETE, OPTIONS",
		AllowHeaders: "content-type",
	}
}

func GetDefaultHeaders() map[string]string {
	return map[string]string{
		"content-type": "text/plain",
		"server":       "Homebrew HTTP Server",
		"connection":   "close",
	}
}
