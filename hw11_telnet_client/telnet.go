package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"time"
)

var ErrNotConnected = errors.New("not connected")

type TelnetClient interface {
	Connect() error
	io.Closer
	Send() error
	Receive() error
}

type telnetClient struct {
	address string
	timeout time.Duration
	in      io.ReadCloser
	out     io.Writer

	connection net.Conn
	scannerIn  *bufio.Scanner
	scannerOut *bufio.Scanner
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &telnetClient{
		address: address,
		timeout: timeout,
		in:      in,
		out:     out,
	}
}

func (c *telnetClient) Connect() error {
	conn, err := net.DialTimeout("tcp", c.address, c.timeout)
	if err != nil {
		return fmt.Errorf("failed to connect to %s: %w", c.address, err)
	}

	c.connection = conn
	c.scannerIn = bufio.NewScanner(c.in)
	c.scannerOut = bufio.NewScanner(c.connection)

	_, err = fmt.Fprintf(os.Stderr, "...Connected to %s\n", c.address)
	if err != nil {
		return err
	}

	return nil
}

func (c *telnetClient) Send() error {
	if c.connection == nil {
		return ErrNotConnected
	}

	if !c.scannerIn.Scan() {
		return io.EOF
	}

	message := fmt.Sprintf("%s\n", c.scannerIn.Text())
	if _, err := c.connection.Write([]byte(message)); err != nil {
		return fmt.Errorf("error sending data: %w", err)
	}

	return nil
}

func (c *telnetClient) Receive() error {
	if c.connection == nil {
		return ErrNotConnected
	}

	if !c.scannerOut.Scan() {
		return io.EOF
	}

	if _, err := fmt.Fprintln(c.out, c.scannerOut.Text()); err != nil {
		return fmt.Errorf("error receiving data: %w", err)
	}

	return nil
}

func (c *telnetClient) Close() error {
	if c.connection == nil {
		return ErrNotConnected
	}

	return c.connection.Close()
}
