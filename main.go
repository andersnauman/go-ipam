package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/andersnauman/go-ipam/api"
	"github.com/andersnauman/go-ipam/config"
	"github.com/andersnauman/go-ipam/ui"
	"github.com/go-sql-driver/mysql"
)

type MUX struct {
	DBpool *sql.DB
}

func (m MUX) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var service http.Handler

	switch strings.SplitN(r.URL.Path, "/", 3)[1] {
	case "api":
		service = api.Object{DatabasePool: m.DBpool}
	default:
		service = ui.Object{DatabasePool: m.DBpool}
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
		User:                 settings.SQL.User,
		Passwd:               settings.SQL.Passwd,
		Net:                  settings.SQL.Net,
		Addr:                 settings.SQL.Addr,
		DBName:               settings.SQL.DBName,
		AllowNativePasswords: settings.SQL.AllowNativePasswords,
		ParseTime:            settings.SQL.ParseTime,
	}

	mux := &MUX{}
	var err error
	mux.DBpool, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal("noooo")
	}

	for {
		pingErr := mux.DBpool.Ping()
		if pingErr == nil {
			break
		} else {
			log.Printf("%s", pingErr.Error())
			time.Sleep(2 * time.Second)
		}
	}

	http.Handle("/", mux)
	http.HandleFunc("/favicon.ico", returnFavicon)
	addr := fmt.Sprintf("%s:%d", settings.Webserver.IP, settings.Webserver.Port)
	log.Printf("%s", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
