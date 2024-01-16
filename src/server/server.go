package server

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"

	"github.com/KSSidll/gohat/src/client"
	"github.com/KSSidll/gohat/src/database"
	"github.com/KSSidll/gohat/src/server/routing"
)

const dbName = "database.db"
const dbPath = dbName

func Serve(port string) {
	db, err := database.ConnectSQLite(dbPath)
	if err != nil {
		log.Fatal(err)
	}

	htmxJs, err := os.ReadFile("src/assets/js/htmx.min.js")
	if err != nil {
		log.Fatal(err)
	}

	mainCSS, err := os.ReadFile("src/assets/css/main.min.css")
	if err != nil {
		log.Fatal(err)
	}

	router := routing.NewRouter()

	// Assets Endpoints
	router.GET("/assets/js/htmx.min.js", func(w *routing.ResponseWriter, r *http.Request) {
		w.ResponseWriter.Header().Add("Content-Type", "text/javascript")
		w.Write(htmxJs)
	}, false)

	router.GET("/assets/css/main.min.css", func(w *routing.ResponseWriter, r *http.Request) {
		w.ResponseWriter.Header().Add("Content-Type", "text/css")
		w.Write(mainCSS)
	}, false)

	// Frontend Endpoints
	router.GET("/", func(w *routing.ResponseWriter, r *http.Request) {
		w.ResponseWriter.Header().Add("Content-Type", "text/html")

		client.ServeIndex(w, r)
	}, false)

	router.GET("/book/all", func(w *routing.ResponseWriter, r *http.Request) {
		w.ResponseWriter.Header().Add("Content-Type", "text/html")

		err := client.ServeBookAllTableComponent(w, r, db)
		if err != nil {
			log.Fatal(err)
		}
	}, false)

	// Backend Endpoints
	router.POST("/book", func(w *routing.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Failed to parse form data", http.StatusBadRequest)
			return
		}

		bookName := r.Form.Get("book-name")

		db.InsertBook(database.BookInsert{Name: bookName})

		w.ResponseWriter.Header().Add("HX-Trigger", "newBook")
	}, true)

	router.DELETE("/book/:id", func(w *routing.ResponseWriter, r *http.Request) {
		id_str := r.Context().Value(routing.ContextKey("id")).(string)

		id, err := strconv.ParseInt(id_str, 10, 64)
		if err != nil {
			log.Fatal(err)
		}

		db.DeleteBook(database.BookDelete{ID: id})

		w.ResponseWriter.Header().Add("HX-Trigger", "deleteBook")
	}, true)

	router.PUT("/book/:id", func(w *routing.ResponseWriter, r *http.Request) {
		id_str := r.Context().Value(routing.ContextKey("id")).(string)

		id, err := strconv.ParseInt(id_str, 10, 64)
		if err != nil {
			log.Fatal(err)
		}

		if err := r.ParseForm(); err != nil {
			http.Error(w, "Failed to parse form data", http.StatusBadRequest)
			return
		}

		bookName := r.Form.Get("book-name")

		db.UpdateBook(database.Book{ID: id, Name: bookName})

		w.ResponseWriter.Header().Add("HX-Trigger", "updateBook")
	}, true)

	l, err := net.Listen("tcp", ":"+port)

	if err != nil {
		fmt.Printf("Error starting server: %s\n", err)
	}

	fmt.Println("Server started on", l.Addr().String())

	if err := http.Serve(l, router); err != nil {
		fmt.Printf("Server closed: %s\n", err)
	}
}
