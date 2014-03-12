package main

import (
	"flag"
	"fmt"
	"os"
	"time"
)

var (
	// Number of questions
	numQuestions = flag.Uint("questions", 10, "")

	// Number of contestants
	numContestants = flag.Uint("contestants", 2, "")

	// Timeout for each question
	timeout = flag.Duration("timeout", 5*time.Second, "")

	// Port to bind
	port = flag.Uint("port", 8888, "")

	// Displays usage
	help = flag.Bool("h", false, "")
)

func main() {
	// Parses the parameters
	flag.Usage = usage
	flag.Parse()

	if *help {
		flag.Usage()
		os.Exit(0)
	}

	server := NewServer()
	server.Listen(*port)
}

func usage() {
	fmt.Println(`Math racer server.

Usage:
	server [arguments]

Arguments:
	-h              Display this help and exit
	--port=8888     Port to bind
	--questions=10  Number of questions to ask to contestants
	--contestants=2 Number of contestants to accept before race is started
	--timeout=5     Timeout for each question
`)
}
