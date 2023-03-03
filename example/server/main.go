package main

import (
	"context"
	"fmt"
	"github.com/andrewwormald/fileuploader"
	"io"
)

func main() {
	ctx := context.Background()
	err := fileuploader.Receive(ctx, func(chunk io.Reader, size int64) error {
		b, err := io.ReadAll(chunk)
		if err != nil {
			return err
		}

		fmt.Println(string(b), size)

		return nil
	}, 8080)
	if err != nil {
		panic(err)
	}
}
