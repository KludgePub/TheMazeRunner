package api

import (
	"bufio"
	"log"
	"net"
	"strings"
)

const (
	STOP string = "STOP_SERVER"
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

func (s *ServerApi) Handle() error {
	defer s.l.Close()

	c, err := s.l.Accept()
	if err != nil {
		return err
	}

	log.Printf("-> Accepted new connection from %s...", c.RemoteAddr())

	// TODO Feed with game map
	c.Write(s.m)

	for {
		netData, err := bufio.NewReader(c).ReadString('\n')
		if err != nil {
			return err
		}

		// Signals
		// TODO Refactor to signal to shutdown server
		if strings.TrimSpace(string(netData)) == STOP {
			log.Println("-> Shutdown game server...")
			return nil
		}

		//fmt.Print("-> ", string(netData))
		//t := time.Now()
		//myTime := t.Format(time.RFC3339) + "\n"
		//c.Write([]byte(myTime))
	}
}
