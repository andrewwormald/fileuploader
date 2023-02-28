package fileuploader

import (
	"context"
	"encoding/binary"
	"fmt"
	"github.com/luno/jettison/log"
	"io"
	"net"
)

func Receive(ctx context.Context, handler func(chunk io.Reader, size int64) error, port int) error {
	ls, err := net.Listen("tcp", fmt.Sprintf(":%v", port))
	if err != nil {
		return err
	}
	defer ls.Close()

	for {
		// Gracefully exit on cancellation of context
		if ctx.Err() != nil {
			return nil
		}

		conn, err := ls.Accept()
		if err != nil {
			return err
		}

		go read(ctx, conn, handler)
	}
}

func read(ctx context.Context, conn net.Conn, handler func(chunk io.Reader, size int64) error) {
	for {
		// Gracefully exit on cancellation of context
		if ctx.Err() != nil {
			return
		}

		var size int64
		err := binary.Read(conn, binary.LittleEndian, &size)
		if err != nil {
			log.Error(ctx, err)
			return
		}

		err = handler(conn, size)
		if err != nil {
			log.Error(ctx, err)
			return
		}
	}
}
