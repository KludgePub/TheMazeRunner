package api

import (
	"bufio"
	"log"
	"net"
	"strings"
)

const (
	stopServer  string = "STOP_SERVER"
	getJsonMaze        = "GET_JSON_MAZE_MAP"
)

// ServerApi used to communicate with game client
type ServerApi struct {
	p string
	l net.Listener
	m []byte // serialized maze graph in json
}

// NewServer create new tcp server api
func NewServer(port string, jMazeMap []byte) (s *ServerApi, err error) {
	s = &ServerApi{p: port, m: jMazeMap}

	s.l, err = net.Listen("tcp", ":"+s.p)
	if err != nil {
		return nil, err
	}

	return s, nil
}

// Handle TCP requests from clients
func (s *ServerApi) Handle() error {
	defer s.l.Close()

	c, err := s.l.Accept()
	if err != nil {
		return err
	}

	log.Printf("-> Accepted new connection from %s...", c.RemoteAddr())

	for {
		netData, err := bufio.NewReader(c).ReadString('\n')
		if err != nil {
			return err
		}
		sReq := string(netData)

		// Handle tcp requests
		if strings.TrimSpace(sReq) == getJsonMaze {
			log.Println("-> Client requested maze map...")
			_, err = c.Write(s.m)
			if err != nil {
				return err
			}
		}

		if strings.TrimSpace(sReq) == stopServer {
			log.Println("-> Client requested to shutdown game server...")
			return nil
		}
	}
}
