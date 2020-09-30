package main

import (
	"encoding/json"
	"net/http"
	"tp-ISA-go.org/kinetur/pkg/forms"
	"tp-ISA-go.org/kinetur/pkg/models"
	_ "tp-ISA-go.org/kinetur/pkg/models"
)

func (app *application) healthCheck(w http.ResponseWriter, _ *http.Request) {
	// Handler de ejemplo que devuele un Json indicando que el servidor esta ok

	// Creo un struct anonima con los valores que quiero mandar
	response := struct {
		Key   string
		Value string
	}{
		"servidor",
		"ok",
	}

	// convierto la string en un json
	js, err := json.Marshal(response)
	if err != nil {
		app.serverError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(js)
	if err != nil {
		app.serverError(w, err)
		return
	}
}

func (app *application) home(w http.ResponseWriter, r *http.Request) {

	app.render(w, r, "home.page.tmpl", nil)
}

func (app *application) signupUserForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "signup.page.tmpl", &templateData{
		Form: forms.New(nil),
	})

}

func (app *application) signupUser(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	//validacion de los datos de usuario
	form := forms.New(r.PostForm)
	form.Required("tipo", "nombre", "apellido", "dni", "domicilio", "email", "password")
	form.MatchesPattern("email", forms.EmailRX)
	form.MinLength("dni", 8)
	form.MinLength("password", 8)

	// si hay un error vuelvo a mostrar el formulario
	if !form.Valid() {
		app.render(w, r, "signup.page.tmpl", &templateData{Form: form})
		return
	}
	// Try to create a new user record in the database. If the email already exi
	// add an error message to the form and re-display it.
	err = app.users.Insert(form.Get("tipo"), form.Get("nombre"), form.Get("apellido"), form.Get("dni"), form.Get("domicilio"), form.Get("email"), form.Get("password"))
	if err == models.ErrDuplicateEmail {
		form.Errors.Add("email", "Direccion ya registrada")
		app.render(w, r, "signup.page.tmpl", &templateData{Form: form})
		return
	} else if err != nil {
		app.serverError(w, err)
		return
	}
	// Otherwise add a confirmation flash message to the session confirming tha
	// their signup worked and asking them to log in.
	app.session.Put(r, "flash", "Registro exitoso, inicie sesion.")

	// And redirect the user to the login page.
	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

func (app *application) loginUserForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "login.page.tmpl", &templateData{
		Form: forms.New(nil),
	})
}
func (app *application) loginUser(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	// Check whether the credentials are valid. If they're not, add a generic e
	// message to the form failures map and re-display the login page.
	form := forms.New(r.PostForm)
	id, err := app.users.Authenticate(form.Get("email"), form.Get("password"))
	if err == models.ErrInvalidCredentials {
		form.Errors.Add("generic", "Email o Password incorrectos")
		app.render(w, r, "login.page.tmpl", &templateData{Form: form})
		return
	} else if err != nil {
		app.serverError(w, err)
		return
	}
	// Add the ID of the current user to the session, so that they are now 'logged in'.
	app.session.Put(r, "userID", id)
	// Redireccion a ver turnos.
	http.Redirect(w, r, "/user/turn", http.StatusSeeOther)
}
func (app *application) logoutUser(w http.ResponseWriter, r *http.Request) {
	// Remove the userID from the session data so that the user is 'logged out'
	app.session.Remove(r, "userID")
	// Add a flash message to the session to confirm to the user that they've be
	app.session.Put(r, "flash", "Ha cerrado sesion con exito")
	http.Redirect(w, r, "/", 303)
}

func (app *application) turnoList(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "turn.page.tmpl", nil)

}

//func onSignIn(googleUser) {
//const googleJWT = googleUser.getAuthResponse().id_token
//}
