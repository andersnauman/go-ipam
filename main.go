package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/go-sql-driver/mysql"

	"github.com/andersnauman/go-ipam/api"
	"github.com/andersnauman/go-ipam/config"
	"github.com/andersnauman/go-ipam/ui"
)

type MUX struct {
	DBpool *sql.DB
}

func (m MUX) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var service http.Handler

	switch strings.SplitN(r.URL.Path, "/", 3)[1] {
	case "api":
		service = api.Object{DatabasePool: nil}
	default:
		service = ui.Object{DatabasePool: nil}
	}

	service.ServeHTTP(w, r)
}

func returnFavicon(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./ui/static/icons/favicon.png")
}

func main() {
	settings := config.Settings

	fs := http.FileServer(http.Dir("./ui/static"))
	http.Handle("/static/", http.StripPrefix("/static", fs))

	cfg := mysql.Config{
		User:   settings.SQL.User,
		Passwd: settings.SQL.Passwd,
		Net:    settings.SQL.Net,
		Addr:   settings.SQL.Addr,
		DBName: settings.SQL.DBName,
	}

	mux := &MUX{}
	var err error
	mux.DBpool, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal("noooo")
	}
	pingErr := mux.DBpool.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}

	rows, err := mux.DBpool.Query("SELECT * FROM album WHERE artist = ?", "John Coltrane")
	if err != nil {
		log.Fatal(fmt.Errorf("albumsByArtist %q: %v", "John Coltrane", err))
		return
	}
	defer rows.Close()
	type Album struct {
		ID     int64
		Title  string
		Artist string
		Price  float32
	}
	for rows.Next() {
		var alb Album
		if err := rows.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price); err != nil {
			log.Fatal(fmt.Errorf("albumsByArtist %q: %v", "John Coltrane", err))
			return
		}
		fmt.Printf(alb.Title)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(fmt.Errorf("albumsByArtist %q: %v", "John Coltrane", err))
		return
	}

	http.Handle("/", mux)
	http.HandleFunc("/favicon.ico", returnFavicon)
	addr := fmt.Sprintf("%s:%d", settings.Webserver.IP.String(), settings.Webserver.Port)
	log.Fatal(http.ListenAndServe(addr, nil))
}
