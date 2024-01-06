package server

import (
	"fmt"
	"net"
	"net/http"

	"github.com/KSSidll/gohat/server/routing"
)

func Serve() {
	router := routing.NewRouter()

	router.GET("/", func(w *routing.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "test")
	})

	l, err := net.Listen("tcp", ":8080")

	if err != nil {
		fmt.Printf("Error starting server: %s\n", err)
	}

	fmt.Println("Server started on", l.Addr().String())

	if err := http.Serve(l, router); err != nil {
		fmt.Printf("Server closed: %s\n", err)
	}
}
