package main

import (
	"fmt"

	"github.com/davidp-ro/homebrew-http-server/tcp"
)

func main() {
	tcp.Server{
		Network: tcp.NetworkTCP,
		Address: ":8080",
		OnConnection: func(data []byte, err error) {
			if err != nil {
				panic(err)
			}

			fmt.Printf("Received: %s\n", data)
		},
	}.Start()
}
