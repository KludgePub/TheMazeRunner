package api

import (
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
	l net.PacketConn
	m []byte // serialized maze graph in json
}

// NewServerConnection create new tcp server api
func NewServerConnection(port string, jMazeMap []byte) (s *ServerApi, err error) {
	s = &ServerApi{p: port, m: jMazeMap}

	s.l, err = net.ListenPacket("udp", "127.0.0.1:"+s.p)
	if err != nil {
		return nil, err
	}

	return s, nil
}

// Handle TCP requests from clients
func (s *ServerApi) Handle() (isClosed bool, handleErr error) {
	defer s.l.Close()

	log.Printf("%s%s\n", "-> UDP API server now handled at ", s.l.LocalAddr())
	for {
		buffer := make([]byte, 1024)
		_, from, readErr := s.l.ReadFrom(buffer)
		if readErr != nil {
			log.Printf("-> UDP API server, incoming read error: %v", readErr.Error())
			continue
		}

		netData := string(buffer)
		log.Printf("-> UDP API server, incoming data (%s) from: %s", netData, from)

		// Handle tcp requests
		if strings.Contains(netData, getJsonMaze) {
			log.Println("-> UDP API server, client requested maze map...")
			_, writeErr := s.l.WriteTo(s.m, from)
			if writeErr != nil {
				log.Printf("-> TCP API server, writing response error: %v", writeErr.Error())
				continue
			}
		} else if strings.Contains(netData, stopServer) {
			log.Println("-> UDP API server, client requested to shutdown game server...")
			return true, nil
		}
	}
}
