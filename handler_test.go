package fileuploader_test

import (
	"bytes"
	"context"
	"encoding/binary"
	"io"
	"net"
	"strings"
	"sync"
	"testing"

	"github.com/luno/jettison/jtest"
	"github.com/stretchr/testify/require"

	"github.com/andrewwormald/fileuploader"
)

func TestReceive(t *testing.T) {
	t.Run("Receive and read chunks of data from tcp connection", func(t *testing.T) {
		var wg sync.WaitGroup
		ctx := context.Background()
		buf := bytes.NewBuffer([]byte{})
		go func() {
			err := fileuploader.Receive(ctx, func(chunk io.Reader, size int64) error {
				// For testing write all chunks to a buffer for examination at the end
				_, err := io.CopyN(buf, chunk, size)
				jtest.RequireNil(t, err)

				wg.Done()
				return nil
			}, 5500)
			jtest.RequireNil(t, err)
		}()

		conn, err := net.Dial("tcp", ":5500")
		jtest.RequireNil(t, err)
		defer conn.Close()

		data := []string{
			"I am the data",
			" -",
			" Hello",
			" World",
		}

		wg.Add(len(data))

		// Mimic a client streaming chunks. A client is expected to provide chunk size then the chunk.
		for _, d := range data {
			b := []byte(d)

			err := binary.Write(conn, binary.LittleEndian, int64(len(b)))
			jtest.RequireNil(t, err)

			_, err = conn.Write(b)
			jtest.RequireNil(t, err)
		}

		wg.Wait()
		received, err := io.ReadAll(buf)
		jtest.RequireNil(t, err)
		require.Equal(t, strings.Join(data, ""), string(received))
	})
}
