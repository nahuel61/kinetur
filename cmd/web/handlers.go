package main

import (
	"net/http"
	"tp-ISA-go.org/kinetur/pkg/forms"
	"tp-ISA-go.org/kinetur/pkg/models"
	_ "tp-ISA-go.org/kinetur/pkg/models"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {

	app.render(w, r, "home.page.tmpl", nil)
}

/*
func (app *application) crearUsuario(w http.ResponseWriter, r *http.Request) {

	nombre := "juan"
	apellido := "Perez"
	dni := "32425219"

	id, err := app.users.Insert(nombre, apellido, dni)
	if err != nil {
		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/user/%d",id), http.StatusSeeOther)

	id, err = strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	//Include the footer partial in the template files.
	app.render(w, r, "signup.page.tmpl", nil)
}

func (app *application) iniciarSesion(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/login" {
		app.notFound(w)
		return
	}
	//Include the footer partial in the template files.

	app.render(w, r, "login.page.tmpl", nil)
}

func (app *application) mostrarUsuario(w http.ResponseWriter, r *http.Request) {

	id, err := strconv.Atoi(r.URL.Query().Get(":id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	_, err = app.users.Get(id)
	if err == models.ErrNoRecord {
		app.notFound(w)
		return
	} else if err != nil {
		app.serverError(w, err)
		return
	}

	// Use the new render helper.
	app.render(w, r, "show.page.tmpl" , &templateData{Usuarios: nil })
}

func(app *application) crearUsuarioForm( w http.ResponseWriter, r *http.Request){
	w.Write([]byte("Crear un nuevo usuario..."))
}
*/

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
		form.Errors.Add("email", "Address is already in use")
		app.render(w, r, "signup.page.tmpl", &templateData{Form: form})
		return
	} else if err != nil {
		app.serverError(w, err)
		return
	}
	// Otherwise add a confirmation flash message to the session confirming tha
	// their signup worked and asking them to log in.
	app.session.Put(r, "flash", "Your signup was successful. Please log in.")

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
	} // Check whether the credentials are valid. If they're not, add a generic e
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
	} // Add the ID of the current user to the session, so that they are now 'logg
	// in'.
	app.session.Put(r, "userID", id)
	// Redirect the user to the create snippet page.
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
func (app *application) logoutUser(w http.ResponseWriter, r *http.Request) {
	// Remove the userID from the session data so that the user is 'logged out'
	app.session.Remove(r, "userID")
	// Add a flash message to the session to confirm to the user that they've be
	app.session.Put(r, "flash", "You've been logged out successfully!")
	http.Redirect(w, r, "/", 303)
}
