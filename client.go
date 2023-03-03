package fileuploader

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
)

func NewClient(address int) *client {
	return &client{
		address: address,
	}
}

type client struct {
	address int
}

func (c *client) Send(r io.Reader, chunks int) error {
	conn, err := net.Dial("tcp", fmt.Sprintf(":%v", c.address))
	if err != nil {
		return err
	}
	defer conn.Close()

	data, err := io.ReadAll(r)
	if err != nil {
		return err
	}

	chunkSize := int64(len(data) / chunks)

	for i := 1; i <= chunks; i++ {
		start := int(chunkSize) * (i - 1)
		end := int(chunkSize) * i
		chunk := data[start:end]

		err = binary.Write(conn, binary.LittleEndian, int64(len(chunk)))
		if err != nil {
			return err
		}

		_, err = conn.Write(chunk)
		if err != nil {
			return err
		}
	}

	return nil
}
