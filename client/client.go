package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
	"strconv"
)

var (
	// Server address
	serverAddr = flag.String("server_address", "0.0.0.0", "")

	// Chat server port
	serverPort = flag.Uint("server_port", 8888, "")

	// Displays usage
	help = flag.Bool("h", false, "")
)

func main() {
	// Parse parameters
	flag.Usage = usage
	flag.Parse()

	if *help {
		flag.Usage()
		os.Exit(0)
	}

	conn, err := net.Dial("tcp", *serverAddr+":"+strconv.Itoa(int(*serverPort)))
	defer conn.Close()
	checkError(err)

	inputChan := make(chan string)
	serverInChan := make(chan string)
	serverOutChan := make(chan string)

	go listenInput(inputChan)
	go listenServerIn(conn, serverInChan)
	go listenServerOut(conn, serverOutChan)

	for {
		select {
		case msg := <-inputChan: // Got input from client's user and forward it to server
			serverOutChan <- msg
		case msgIn := <-serverInChan: // Got new msg from server, print it.
			// We expect msgIn delivered with '\n'.
			fmt.Print(msgIn)
		}
	}

}

// listenInput listens for stdin and send back
// to input channel.
func listenInput(inputchan chan<- string) {
	input := bufio.NewReader(os.Stdin)

	for {
		line, _, err := input.ReadLine()
		checkError(err)
		inputchan <- fmt.Sprintf("%s\n", string(line))
	}
}

// listenServerIn listens message from connection
// and forward it to serverInChan
func listenServerIn(conn net.Conn, serverInChan chan string) {
	chatReader := bufio.NewReader(conn)

	for {
		msgIn, err := chatReader.ReadString('\n')
		checkError(err)

		serverInChan <- string(msgIn)
	}
}

// listenServerOut listens on serverOutChan channel and forward
// the message to server connection.
func listenServerOut(conn net.Conn, serverOutChan <-chan string) {
	for {
		msg := <-serverOutChan
		conn.Write([]byte(msg))
	}
}

// usages shows usage.
func usage() {
	fmt.Println(`Chat client.

Usage:
	client [arguments]

Arguments:
	-h                         Display this help and exit
	--server_port=8888         Chat server port
	--server_address="0.0.0.0" Chat server address
`)
}

// checkError checks error, will exit if error is not nil.
func checkError(err error) {
	if err != nil {
		fmt.Printf("Error: %v\n", err.Error())
		os.Exit(1)
	}
}
