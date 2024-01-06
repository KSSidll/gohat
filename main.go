package main

import (
	"os"

	"github.com/KSSidll/gohat/server"
)

func main() {
	server.Serve()
	os.Exit(1)
}
