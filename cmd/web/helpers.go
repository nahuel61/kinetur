package main

import (
	"bytes"
	"fmt"
	"github.com/justinas/nosurf"
	"net/http"
	"runtime/debug"
	"time"
	"tp-ISA-go.org/kinetur/pkg/models"
)

// 500 Internal Server Error response to the user en el log
func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf(" %s\n%s", err.Error(), debug.Stack())
	app.errorLog.Output(2, trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// Me da las respuestas que veo en el postman (ie: 400 Bad request)
func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

// 404 not found
func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}

func (app *application) addDefaultData(td *templateData, r *http.Request) *templateData {
	if td == nil {
		td = &templateData{}
	}
	td.CSRFToken = nosurf.Token(r)
	td.AuthenticatedUser = app.authenticatedUser(r)
	td.Año = time.Now().Year()
	td.Flash = app.session.PopString(r, "flash") //agrego el mensaje a template data, si existe lo muetra.
	return td
}

//funcion que renderiza los templates html
func (app *application) render(w http.ResponseWriter, r *http.Request, name string, td *templateData) {
	ts, ok := app.templateCache[name]
	if !ok {
		app.serverError(w, fmt.Errorf("El template %s NO existe", name))
		return
	}
	// Escribo el template en un buffer que me permite
	//mostrar si hay un error antes de llamar al http.ResponseWriter.
	buf := new(bytes.Buffer) //creo un nuevo buffer
	err := ts.Execute(buf, app.addDefaultData(td, r))
	if err != nil {
		app.serverError(w, err)
		return
	}
	buf.WriteTo(w) // paso lo del buffer al http.ResponseWriter

}

// si el usuario esta autenticado va a responder con la contextkeyuser
func (app *application) authenticatedUser(r *http.Request) *models.Pacientes {
	user, ok := r.Context().Value(contextKeyUser).(*models.Pacientes)
	if !ok {
		return nil
	}
	return user
}
