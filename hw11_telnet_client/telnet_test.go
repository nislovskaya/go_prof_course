package main

import (
	"bytes"
	"io"
	"net"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestTelnetClient(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		l, err := net.Listen("tcp", "127.0.0.1:")
		require.NoError(t, err)
		defer func() { require.NoError(t, l.Close()) }()

		var wg sync.WaitGroup
		wg.Add(2)

		go func() {
			defer wg.Done()

			in := &bytes.Buffer{}
			out := &bytes.Buffer{}

			timeout, err := time.ParseDuration("10s")
			require.NoError(t, err)

			client := NewTelnetClient(l.Addr().String(), timeout, io.NopCloser(in), out)
			require.NoError(t, client.Connect())
			defer func() { require.NoError(t, client.Close()) }()

			in.WriteString("hello\n")
			err = client.Send()
			require.NoError(t, err)

			err = client.Receive()
			require.NoError(t, err)
			require.Equal(t, "world\n", out.String())
		}()

		go func() {
			defer wg.Done()

			conn, err := l.Accept()
			require.NoError(t, err)
			require.NotNil(t, conn)
			defer func() { require.NoError(t, conn.Close()) }()

			request := make([]byte, 1024)
			n, err := conn.Read(request)
			require.NoError(t, err)
			require.Equal(t, "hello\n", string(request)[:n])

			n, err = conn.Write([]byte("world\n"))
			require.NoError(t, err)
			require.NotEqual(t, 0, n)
		}()

		wg.Wait()
	})

	t.Run("connection timeout", func(t *testing.T) {
		client := NewTelnetClient("127.0.0.1:12345", 1*time.Second, io.NopCloser(bytes.NewBufferString("test\n")), io.Discard)
		err := client.Connect()
		require.Error(t, err)
		require.Contains(t, err.Error(), "failed to connect to 127.0.0.1:12345")
	})

	t.Run("send after close", func(t *testing.T) {
		l, err := net.Listen("tcp", "127.0.0.1:")
		require.NoError(t, err)
		defer func() { require.NoError(t, l.Close()) }()

		var wg sync.WaitGroup
		wg.Add(1)

		go func() {
			defer wg.Done()

			conn, err := l.Accept()
			require.NoError(t, err)
			defer func() { require.NoError(t, conn.Close()) }()

			time.Sleep(100 * time.Millisecond)
		}()

		in := &bytes.Buffer{}
		out := &bytes.Buffer{}
		timeout := 10 * time.Second

		client := NewTelnetClient(l.Addr().String(), timeout, io.NopCloser(in), out)
		require.NoError(t, client.Connect())

		require.NoError(t, client.Close())

		in.WriteString("test\n")
		err = client.Send()
		require.ErrorContains(t, err, "error sending data")

		wg.Wait()
	})

	t.Run("no server", func(t *testing.T) {
		client := NewTelnetClient(
			"localhost:4242",
			time.Second*10,
			io.NopCloser(&bytes.Buffer{}),
			&bytes.Buffer{},
		)

		err := client.Connect()
		require.Error(t, err)
		require.Contains(t, err.Error(), "failed to connect to localhost:4242")
	})

	t.Run("not connected", func(t *testing.T) {
		client := NewTelnetClient(
			"localhost:4242",
			time.Second*10,
			io.NopCloser(&bytes.Buffer{}),
			&bytes.Buffer{},
		)

		sendErr := client.Send()
		require.ErrorIs(t, sendErr, ErrNotConnected)

		receiveErr := client.Receive()
		require.ErrorIs(t, receiveErr, ErrNotConnected)

		closeErr := client.Close()
		require.ErrorIs(t, closeErr, ErrNotConnected)
	})
}
