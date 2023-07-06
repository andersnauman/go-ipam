package ui

import (
	"database/sql"
	"html"
	"log"
	"net/http"
	"strings"
	"text/template"

	"github.com/andersnauman/go-ipam/misc"
)

type Object struct {
	DatabasePool *sql.DB
}

type PageIndex struct {
	PageName string
	User     misc.User
}

type PageDevices struct {
	PageName string
	User     misc.User
}

var templates = template.Must(template.ParseFiles(
	"./ui/templates/_header.tmpl",
	"./ui/templates/_footer.tmpl",
	"./ui/templates/index.tmpl",
	"./ui/templates/devices.tmpl",
	"./ui/templates/locations.tmpl",
	"./ui/templates/login.tmpl",
	"./ui/templates/modals/modal_devices_edit.tmpl",
))

func (o Object) rendertemplate(w http.ResponseWriter, tmpl string, t interface{}) {
	err := templates.ExecuteTemplate(w, tmpl+".tmpl", t)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (o Object) pageDevices(w http.ResponseWriter, r *http.Request, path []string) (err error) {
	o.rendertemplate(w, "devices", PageDevices{
		PageName: "devices",
		User: misc.User{
			FirstName: "test",
			LastName:  "testsson",
		},
	})
	return
}

func (o Object) modalDevicesEdit(w http.ResponseWriter, r *http.Request, path []string) (err error) {
	o.rendertemplate(w, "modal_devices_edit", nil)
	return
}

func (o Object) pageLocations(w http.ResponseWriter, r *http.Request, path []string) (err error) {
	var authorized bool
	authorized, err = misc.VerifyLogin(r, o.DatabasePool)
	if !authorized {
		http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
	}
	o.rendertemplate(w, "locations", map[string]interface{}{
		"PageName": "locations",
		"User": misc.User{
			FirstName: "test",
			LastName:  "testsson",
		},
		"Locations": []map[string]interface{}{{
			"Coordinates": "59.334591, 18.063240",
			"Description": "Datahall 1<br>1.1.1.0/24",
		}, {
			"Coordinates": "59.334691, 18.063240",
			"Description": "Datahall 2<br>1.1.2.0/24",
		},
		}})
	return
}

func (o Object) pageLogin(w http.ResponseWriter, r *http.Request, path []string) (err error) {
	var authorized bool
	authorized, err = misc.VerifyLogin(r, o.DatabasePool)
	if !authorized {
		o.rendertemplate(w, "login", nil)
		return
	}

	err = misc.AddSessionCookie(w, r, o.DatabasePool)

	http.Redirect(w, r, "/index", http.StatusSeeOther)
	return
}

func (o Object) pageIndex(w http.ResponseWriter, r *http.Request, path []string) (err error) {
	o.rendertemplate(w, "index", PageIndex{
		PageName: "index",
		User: misc.User{
			FirstName: "test",
			LastName:  "testsson",
		},
	})
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
		err = o.pageLogin(w, r, path)
		if err != nil {
			log.Printf("Login error: %s", err.Error())
		}
	case "devices":
		_ = o.pageDevices(w, r, path)
	case "locations":
		_ = o.pageLocations(w, r, path)
	case "modal":
		log.Println(path[1])
		_ = o.modalDevicesEdit(w, r, path)
	default:
		_ = o.pageIndex(w, r, path)
	}
	//http.Error(w, "neee", http.StatusInternalServerError)
}
