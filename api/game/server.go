package game

import (
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/LinMAD/TheMazeRunnerServer/manager"
)

const logTag = "-> UDP API server:"

const (
	stopServer  string = "STOP_SERVER"
	getJsonMaze        = "GET_JSON_MAZE_MAP"
	getPlayerMoves     = "GET_PLAYER_MOVES"
)

// UDPServerAPI used to communicate with game client
type UDPServerAPI struct {
	// forceClose of UDP server
	forceClose chan os.Signal
	p          string
	l          net.PacketConn
	m          []byte // serialized maze graph in json
}

// NewServerConnection create new tcp server api
func NewServerConnection(port string, jMazeMap []byte) (s *UDPServerAPI, err error) {
	s = &UDPServerAPI{p: port, m: jMazeMap, forceClose: make(chan os.Signal, 1)}

	s.l, err = net.ListenPacket("udp", "127.0.0.1:"+s.p)
	if err != nil {
		return nil, err
	}

	return s, nil
}

// Handle TCP requests from clients
func (s *UDPServerAPI) Handle(gm *manager.GameManager) (isClosed bool, handleErr error) {
	defer s.l.Close()

	go func(isClosed bool) {
		signal.Notify(s.forceClose, os.Interrupt, syscall.SIGTERM)
		<-s.forceClose

		log.Printf("%s %s", logTag, "Performing graceful shutdown of UDP API...")

		isClosed = true
	}(isClosed)

	log.Printf("%s %s %s\n", logTag, "now handled at", s.l.LocalAddr())

	for {
		if isClosed == true {
			handleErr = s.l.Close()
			break
		}

		buffer := make([]byte, 1024)
		_, from, readErr := s.l.ReadFrom(buffer)
		if readErr != nil {
			log.Printf("%s incoming read error: %v", logTag, readErr.Error())
			continue
		}

		netData := string(buffer)
		log.Printf("%s incoming data (%s) from: %s", logTag, netData, from)

		// Handle UDP requests
		if strings.Contains(netData, getJsonMaze) {
			log.Printf("%s client requested maze map...\n", logTag)

			_, writeErr := s.l.WriteTo(s.m, from)
			if writeErr != nil {
				log.Printf("%s writing response error: %v\n", logTag, writeErr.Error())
				continue
			}

		} else if strings.Contains(netData, stopServer) {
			log.Printf("%s client requested to shutdown game server...\n", logTag)
			isClosed = true
		} else if strings.Contains(netData, getPlayerMoves) {
			_, writeErr := s.l.WriteTo(gm.GetPlayerMovements(), from)
			if writeErr != nil {
				log.Printf("%s writing response error: %v\n", logTag, writeErr.Error())
				continue
			}
		}
	}

	return isClosed, handleErr
}
