package main

import (
	"bufio"
	"net"
)

// Client represents connected client in chat server.
type Client struct {
	nick   string
	in     chan Message
	out    chan string
	quit   chan string
	reader *bufio.Reader
	writer *bufio.Writer
}

// ReadListener listens any message from tcp connection that
// this client sent and then forward the message to server's
// incoming channel.
func (c *Client) ReadListener() {
	for {
		line, _, err := c.reader.ReadLine()
		if err != nil {
			break
		}

		msg := Message{
			nick:    c.nick,
			message: string(line),
		}
		c.in <- msg
	}

	c.quit <- c.nick
}

// WriteListener writes message to client connection
// from outgoing channel. Outgoing channel will be
// filled by server such as when broadcasting message.
func (c *Client) WriteListener() {
	for msg := range c.out {
		c.writer.WriteString(msg)
		c.writer.Flush()
	}

}

// Listen listen incoming and outgoing message.
func (c *Client) Listen() {
	go c.ReadListener()
	go c.WriteListener()
}

// NewClient returns a new Client that's listening for
// any incoming and outgoing message.
func NewClient(conn net.Conn, nick string) *Client {
	writer := bufio.NewWriter(conn)
	reader := bufio.NewReader(conn)

	client := &Client{
		nick:   nick,
		in:     make(chan Message),
		out:    make(chan string),
		quit:   make(chan string),
		reader: reader,
		writer: writer,
	}

	client.Listen()

	return client
}
