package main

import (
	"github.com/andrewwormald/fileuploader"
	"strings"
)

func main() {
	cl := fileuploader.NewClient(8080)

	r := strings.NewReader(`
		Hello there,

		This message is broken up into shards and streamed over tcp.

		I had a lot of fun writing this.
	`)
	err := cl.Send(r, 10)
	if err != nil {
		panic(err)
	}
}
