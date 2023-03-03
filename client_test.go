package fileuploader_test

import (
	"bytes"
	"context"
	"io"
	"strings"
	"sync"
	"testing"

	"github.com/luno/jettison/jtest"
	"github.com/stretchr/testify/require"

	"github.com/andrewwormald/fileuploader"
)

func TestSend(t *testing.T) {
	t.Run("Stream reader to backend", func(t *testing.T) {
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

		data := "I am the data - Hello World"

		cl := fileuploader.NewClient(5500)
		chunks := 3

		wg.Add(chunks)

		err := cl.Send(strings.NewReader(data), chunks)
		jtest.RequireNil(t, err)

		wg.Wait()
		received, err := io.ReadAll(buf)
		jtest.RequireNil(t, err)
		require.Equal(t, data, string(received))
	})
}
