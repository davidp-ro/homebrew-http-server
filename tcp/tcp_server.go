package tcp

import "net"

type TCPNetwork string

type OnConnectionCallback func([]byte, error) []byte

const (
	NetworkTCP  TCPNetwork = "tcp"
	NetworkTCP4 TCPNetwork = "tcp4"
	NetworkTCP6 TCPNetwork = "tcp6"
)

type Server struct {
	Network      TCPNetwork
	Address      string
	OnConnection OnConnectionCallback
}

func (server Server) Start() error {
	listener, err := net.Listen(string(server.Network), server.Address)
	if err != nil {
		return err
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			return err
		}

		go handleConnection(conn, server.OnConnection)
	}
}

func handleConnection(conn net.Conn, onConn OnConnectionCallback) {
	defer conn.Close()

	buf := make([]byte, 1024)
	_, err := conn.Read(buf)

	if err != nil {
		onConn(nil, err)
		conn.Write([]byte("HTTP/1.1 500 Internal Server Error\r\n\r\n"))
	}

	conn.Write(onConn(buf, nil))
}
