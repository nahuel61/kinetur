package main

import (

	"html/template"
	"net/http"
)

func (app *application) home (w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}
	//Include the footer partial in the template files.
		files := []string{
		"./ui/html/home.page.tmpl",
		"./ui/html/base.layout.tmpl",
		"./ui/html/footer.partial.tmpl",
	}
	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(w,err)
		return
	}
	err = ts.Execute(w, nil)
	if err != nil {
		app.serverError(w,err)
	}
}
func (app *application) crearUsuario (w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/register" {
		app.notFound(w)
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
		app.serverError(w,err)
		return
	}
	err = ts.Execute(w, nil)
	if err != nil {
		app.serverError(w,err)
	}
}
func (app *application) iniciarSesion(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/login" {
		app.notFound(w)
		return
	}
	//Include the footer partial in the template files.
	files := []string{
		"./ui/html/login.page.tmpl",
		"./ui/html/base.layout.tmpl",
		"./ui/html/footer.partial.tmpl",
	}
	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(w,err)
		return
	}
	err = ts.Execute(w, nil)
	if err != nil {
		app.serverError(w,err)
	}
}
