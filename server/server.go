package server

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"strconv"

	"github.com/KSSidll/gohat/client"
	"github.com/KSSidll/gohat/database"
	"github.com/KSSidll/gohat/server/routing"
)

const dbName = "database.db"
const dbPath = dbName

func Serve(port string) {
    db, err := database.ConnectSQLite(dbPath)
    if err != nil {
        log.Fatal(err)
    }

	router := routing.NewRouter()

	router.GET("/", func(w *routing.ResponseWriter, r *http.Request) {
		client.ServeIndex(w, r)
	}, false)

	router.POST("/book", func(w *routing.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Failed to parse form data", http.StatusBadRequest)
			return
		}

		bookName := r.Form.Get("book-name")

		db.InsertBook(database.BookInsert{ Name: bookName })
	}, true)

	router.GET("/book/all", func(w *routing.ResponseWriter, r *http.Request) {
		err := client.ServeBookAllTableComponent(w, r, db)
		if err != nil {
			log.Fatal(err)
		}
	}, false)

	router.DELETE("/book/:id", func(w *routing.ResponseWriter, r *http.Request) {
		id_str :=   r.Context().Value(routing.ContextKey("id")).(string)

		id, err := strconv.ParseInt(id_str, 10, 64)
		if err != nil {
			log.Fatal(err)
		}

		db.DeleteBook(database.BookDelete{ ID: id })
	}, true)

	l, err := net.Listen("tcp", ":" + port)

	if err != nil {
		fmt.Printf("Error starting server: %s\n", err)
	}

	fmt.Println("Server started on", l.Addr().String())

	if err := http.Serve(l, router); err != nil {
		fmt.Printf("Server closed: %s\n", err)
	}
}
