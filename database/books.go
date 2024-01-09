package database

import (
	"errors"

	"github.com/mattn/go-sqlite3"
)

type Book struct {
	ID int64
	Name string
}

type BookInsert struct {
	Name string
}

type BookDelete struct {
	ID int64
}

func (r *SQLiteDatabase) InsertBook(book BookInsert) (int64, error) {
	res, err := r.db.Exec("INSERT INTO books(name) VALUES(?)", book.Name)

    if err != nil {
        var sqliteErr sqlite3.Error
        if errors.As(err, &sqliteErr) {
            if errors.Is(sqliteErr.ExtendedCode, sqlite3.ErrConstraintUnique) {
                return -1, ErrDuplicate
            }
        }
        return -1, err
	}

	id, err := res.LastInsertId()
    if err != nil {
        return -1, err
    }

    return id, nil
}

func (r *SQLiteDatabase) GetBookAll() ([]Book, error) {
	rows, err := r.db.Query("SELECT * FROM books")
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var all []Book
    for rows.Next() {
        var book Book
        if err := rows.Scan(&book.ID, &book.Name); err != nil {
            return nil, err
        }
        all = append(all, book)
    }

    return all, nil
}

func (r *SQLiteDatabase) DeleteBook(book BookDelete) error {
	res, err := r.db.Exec("DELETE FROM books WHERE id = ?", book.ID)
    if err != nil {
        return err
    }

    rowsAffected, err := res.RowsAffected()
    if err != nil {
        return err
    }

    if rowsAffected == 0 {
        return ErrDeleteFailed
    }

    return err

}

func (r *SQLiteDatabase) UpdateBook(book Book) (*Book, error) {
	if book.ID == 0 {
        return nil, errors.New("invalid updated ID")
    }
    res, err := r.db.Exec("UPDATE books SET name = ? WHERE id = ?", book.Name, book.ID)
    if err != nil {
        return nil, err
    }

    rowsAffected, err := res.RowsAffected()
    if err != nil {
        return nil, err
    }

    if rowsAffected == 0 {
        return nil, ErrUpdateFailed
    }

    return &book, nil
}
