package main

import (
	"fmt"

	"github.com/davidp-ro/homebrew-http-server/http"
	"github.com/davidp-ro/homebrew-http-server/tcp"
)

func main() {
	tcp.Server{
		Network: tcp.NetworkTCP,
		Address: ":8080",
		OnConnection: func(data []byte, err error) []byte {
			// fmt.Printf(">> REQUEST:\n%q\n<< END\n\n", data)

			if err != nil {
				panic(err)
			}

			parsed, err := http.ParseRawRequest(data)
			if err != nil {
				panic(err)
			}

			headers := ""
			for k, v := range parsed.Headers {
				headers += fmt.Sprintf("%s: %s\n", k, v)
			}

			fmt.Println(parsed.Body)

			return []byte("HTTP/1.1 200 OK\n\nOK")
			// return []byte(fmt.Sprintf("HTTP/1.1 200 OK\n\n%s\n%s\n%s", parsed.Path, parsed.QueryParams["param"], headers))
		},
	}.Start()
}
