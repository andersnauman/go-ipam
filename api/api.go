package api

import (
	"database/sql"
	"net/http"
)

type Object struct {
	DatabasePool *sql.DB
}

func (o Object) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "jaaa", http.StatusInternalServerError)
}
