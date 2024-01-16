package client

import (
	"net/http"

	"github.com/KSSidll/gohat/src/database"
	"github.com/KSSidll/gohat/src/server/routing"
	"github.com/a-h/templ"
)

func ServeIndex(w *routing.ResponseWriter, r *http.Request) {
	templ.Handler(IndexPage()).ServeHTTP(w, r)
}

func ServeBookAllTableComponent(w *routing.ResponseWriter, r *http.Request, db *database.SQLiteDatabase) error {
	books, err := db.GetBookAll()
	if err != nil {
		return err
	}

	templ.Handler(BookTableComponent(books)).ServeHTTP(w, r)

	return nil
}
