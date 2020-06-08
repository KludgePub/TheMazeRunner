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
func (s *ServerApi) Handle() (isClosed bool, handleErr error) {
	defer s.l.Close()

	c, err := s.l.Accept()
	if err != nil {
		return false, err
	}

	log.Printf("-> TCP API server, accepted new connection from %s...", c.RemoteAddr())

	for {
		netData, err := bufio.NewReader(c).ReadString('\n')
		if err != nil {
			return false, err
		}
		sReq := string(netData)

		// Handle tcp requests
		if strings.TrimSpace(sReq) == getJsonMaze {
			log.Println("-> TCP API server, client requested maze map...")
			_, err = c.Write(s.m)
			if err != nil {
				return false, err
			}
		}

		if strings.TrimSpace(sReq) == stopServer {
			log.Println("-> TCP API server, client requested to shutdown game server...")
			return true, nil
		}
	}
}
