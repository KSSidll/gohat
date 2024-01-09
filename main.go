package main

import (
	"os"

	"github.com/KSSidll/gohat/server"
)

func main() {
	server.Serve("8080")
	os.Exit(1)
}
