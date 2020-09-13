package main

import (
	"fmt"
	_ "fmt"
	"html/template"
	"net/http"
	"strconv"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}

	app.render(w, r, "home.page.tmpl")
}
func (app *application) crearUsuario(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/register" {
		app.notFound(w)
		return
	}

	if r.Method != "POST" {
		return
	}

	nombre := "Nahuel"
	apellido := "Salazar"
	dni := "32424219"

	id, err := app.users.Insert(nombre, apellido, dni)
	if err != nil {
		app.serverError(w, err)
	}
	http.Redirect(w, r, fmt.Sprintf("/"), http.StatusSeeOther)

	id, err = strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	//Include the footer partial in the template files.
	files := []string{
		"./ui/html/register.page.tmpl",
		"./ui/html/base.layout.tmpl",
		"./ui/html/footer.partial.tmpl",
	}
	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(w, err)
		return
	}
	err = ts.Execute(w, nil)
	if err != nil {
		app.serverError(w, err)
	}
}

func (app *application) iniciarSesion(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/login" {
		app.notFound(w)
		return
	}
	//Include the footer partial in the template files.

	app.render(w, r, "login.page.tmpl")
}

func (app *application) mostrarUsuario(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/user" {
		app.notFound(w)
		return
	}

	s, err := app.users.Latest()

	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	data := &templateData{Usuarios: s}

	files := []string{
		"./ui/html/show.page.tmpl",
		"./ui/html/base.layout.tmpl",
		"./ui/html/footer.partial.tmpl",
	}
	ts, err := template.ParseFiles(files...)

	if err != nil {
		app.serverError(w, err)
		return
	}
	err = ts.Execute(w, data)
	if err != nil {
		app.serverError(w, err)
	}

	// Use the new render helper.
	app.render(w, r, "show.page.tmpl")
}
