package main

import (
	"errors"
	"flag"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var logger = log.New(os.Stderr, "telnet: ", log.LstdFlags)

func main() {
	var timeout time.Duration
	flag.DurationVar(&timeout, "timeout", 10*time.Second, "connection timeout")
	flag.Parse()

	if flag.NArg() < 2 {
		logger.Fatal("Usage: go-telnet [--timeout=duration] host port")
	}

	host := flag.Arg(0)
	port := flag.Arg(1)

	sigintChan := make(chan os.Signal, 1)
	signal.Notify(sigintChan, syscall.SIGINT)

	client := NewTelnetClient(net.JoinHostPort(host, port), timeout, os.Stdin, os.Stdout)

	if err := client.Connect(); err != nil {
		logger.Fatal(err)
	}
	defer client.Close()

	go handleReceiving(client)
	go handleSending(client)

	waitForTermination(sigintChan)
}

func handleReceiving(client TelnetClient) {
	for {
		if err := client.Receive(); err != nil {
			logger.Println("Error receiving data:", err)
			break
		}
	}
}

func handleSending(client TelnetClient) {
	for {
		if err := client.Send(); err != nil {
			if errors.Is(err, io.EOF) {
				logger.Println("Input finished. Closing client.")
			} else {
				logger.Println("Error sending data:", err)
			}
			break
		}
	}
}

func waitForTermination(sigintChan chan os.Signal) {
	select {
	case <-sigintChan:
		logger.Println("SIGINT. Closing.")
	case <-time.After(10 * time.Second):
		logger.Println("Wait timeout.")
	}
}
