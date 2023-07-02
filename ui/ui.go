package ui

import (
	"database/sql"
	"fmt"
	"html"
	"net/http"
	"strings"
	"text/template"

	"github.com/andersnauman/go-ipam/misc"
)

type Object struct {
	DatabasePool *sql.DB
}

var templates = template.Must(template.ParseFiles(
	"./ui/templates/_header.tmpl",
	"./ui/templates/_footer.tmpl",
	"./ui/templates/index.tmpl",
	"./ui/templates/devices.tmpl",
	"./ui/templates/locations.tmpl",
	"./ui/templates/login.tmpl",
))

func (o Object) rendertemplate(w http.ResponseWriter, tmpl string, t interface{}) {
	err := templates.ExecuteTemplate(w, tmpl+".tmpl", t)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (o Object) pageDevices(w http.ResponseWriter, r *http.Request, path []string) (err error) {
	o.rendertemplate(w, "devices", nil)
	return
}

func (o Object) pageLocations(w http.ResponseWriter, r *http.Request, path []string) (err error) {
	o.rendertemplate(w, "locations", nil)
	return
}

func (o Object) pageLogin(w http.ResponseWriter, r *http.Request, path []string) (err error) {
	switch r.Method {
	case "POST":
		if err = r.ParseForm(); err != nil {
			fmt.Fprintf(w, "ParseForm() err: %v", err)
			return
		}
		username := r.FormValue("username")
		if len(username) < 1 {
			fmt.Printf("Empty field username")
			return
		}
		password := r.FormValue("password")
		if len(password) < 1 {
			fmt.Printf("Empty field password")
			return
		}
		fmt.Printf("Username: %s, Password: %s", username, password)
		_ = o.pageIndex(w, r, path)
	default:
		o.rendertemplate(w, "login", nil)
	}
	return
}

func (o Object) pageIndex(w http.ResponseWriter, r *http.Request, path []string) (err error) {
	o.rendertemplate(w, "index", nil)
	return
}

func (o Object) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.Split(html.EscapeString(r.URL.Path), "/")
	path, err := misc.PopString(path)
	if err != nil {
		_ = misc.ReturnError(w, http.StatusBadRequest, err.Error())
		return
	}
	site := path[0]
	switch site {
	case "login":
		_ = o.pageLogin(w, r, path)
	case "devices":
		_ = o.pageDevices(w, r, path)
	case "locations":
		_ = o.pageLocations(w, r, path)
	default:
		_ = o.pageIndex(w, r, path)
	}
	//http.Error(w, "neee", http.StatusInternalServerError)
}
