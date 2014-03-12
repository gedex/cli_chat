package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
)

// Server represents chat server
type Server struct {
	clients map[string]*Client
	conn    chan net.Conn
	in      chan Message
	out     chan string
	quit    chan string
}

// Broadcast sends message to all clients except the user
// with nick specified in sender
func (s *Server) Broadcast(sender string, msg string) {
	for nick, c := range s.clients {
		if nick != sender {
			c.out <- fmt.Sprintf("\033[1;33;40m%s>\033[0m \033[1;36m%s\033[0m\n", sender, msg)
		}
	}
}

// NewUserConnection asks new connection a nick to be used
// in chat. Nick must be unique. Once nick is assigned to
// connection a goroutine will be run that listening incoming
// message and quit sign via channel.
func (s *Server) NewUserConnection(conn net.Conn) {
	nick, err := s.askNick(conn)
	if err != nil {
		log.Printf("Connection %v unable to join", conn.RemoteAddr())
		return
	}

	c := NewClient(conn, nick)
	s.clients[nick] = c

	log.Printf(fmt.Sprintf("\033[0;32m%s joined the chat\033[0m", nick))

	// Notify other users about new user joined.
	s.Broadcast(nick, "joined")

	// Welcome the new user.
	conn.Write([]byte(fmt.Sprintf("\033[0;32mWelcome %s\033[0m\n", nick)))

	go func() {
		for {
			select {
			case msg := <-c.in:
				s.in <- msg
			case quit := <-c.quit:
				s.quit <- quit
				conn.Close()
			}
		}
	}()
}

// askNick asks the connection for a nick.
func (s *Server) askNick(conn net.Conn) (string, error) {
	reader := bufio.NewReader(conn)

	// Asks nick until a unique one is satisfied.
	conn.Write([]byte("\033[1;32mPlease enter your nick:\033[0m\n"))
	for {
		nick, _, err := reader.ReadLine()
		if err != nil {
			return "", err
		}

		// Makes sure nick is unique
		if _, exists := s.clients[string(nick)]; exists {
			conn.Write([]byte("\033[1;31mNick already used, please try different one:\033[0m \n"))

			log.Printf("\033[1;31m%v try joining with registered nick '%s'\033[0m", conn.RemoteAddr(), nick)
		} else {
			return string(nick), nil
		}
	}
}

// isNickExists checks whether a given nick already exists.
func (s *Server) isNickExists(nick string) bool {
	if _, exists := s.clients[nick]; exists {
		return true
	}

	return false
}

// NewServer returns a new Server that host the chat server.
func NewServer() *Server {
	server := &Server{
		clients: make(map[string]*Client),
		conn:    make(chan net.Conn),
		in:      make(chan Message),
		out:     make(chan string),
		quit:    make(chan string),
	}

	return server
}

// Listen listens any tcp connection on specified port.
func (s *Server) Listen(port uint) {
	go s.channelReceivingHandler()

	listener, err := net.Listen("tcp", ":"+strconv.Itoa(int(port)))
	if err != nil {
		log.Fatalf("Error: %v\n", err.Error())
		os.Exit(1)
	}
	log.Printf("\033[1;32mServer listening on port %v\033[0m", port)

	for {
		conn, _ := listener.Accept()
		s.conn <- conn
	}
}

// channelReceivingHandler runs infinitely to check
// receing channel (incoming message, new connection, quit)
func (s *Server) channelReceivingHandler() {
	for {
		select {
		case msg := <-s.in:
			s.handleNewMessage(msg)
		case conn := <-s.conn:
			s.handleNewConnection(conn)
		case quit := <-s.quit:
			s.handleUserQuit(quit)
		}
	}
}

// handleNewMessage handles new incoming message by broadcasting
// it to all connected clients.
func (s *Server) handleNewMessage(msg Message) {
	log.Printf("\033[1;33mNew message from %s: '%s'\033[m", msg.nick, msg.message)
	s.Broadcast(msg.nick, msg.message)
}

// handleNewConnection handles a connected client.
func (s *Server) handleNewConnection(conn net.Conn) {
	log.Printf("\033[1;33mNew connection from %v\033[0m", conn.RemoteAddr())
	s.NewUserConnection(conn)
}

// handleUserQuit handles closed connection from particular
// client specified by its nick.
func (s *Server) handleUserQuit(nick string) {
	if _, ok := s.clients[nick]; ok {
		delete(s.clients, nick)

		log.Printf("\033[1;31mUser %s quits\033[0m", nick)
		s.Broadcast(nick, "quit")
	}
}
