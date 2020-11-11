package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"tp-ISA-go.org/kinetur/pkg/forms"
	"tp-ISA-go.org/kinetur/pkg/models"
	_ "tp-ISA-go.org/kinetur/pkg/models"
)

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
	form.MinLength("password", 6)
	// si hay un error vuelvo a mostrar el formulario
	if !form.Valid() {
		app.render(w, r, "signup.page.tmpl", &templateData{Form: form})
		return
	}
	// creo un nuevo registro en la base de datos, si se repite el email avisa
	err = app.users.Insert(form.Get("tipo"), form.Get("nombre"), form.Get("apellido"), form.Get("dni"), form.Get("domicilio"), form.Get("email"), form.Get("password"))
	if err == models.ErrDuplicateEmail {
		form.Errors.Add("email", "Direccion ya registrada")
		app.render(w, r, "signup.page.tmpl", &templateData{Form: form})
		return
	} else if err != nil {
		app.serverError(w, err)
		return
	}
	// aviso que el registro fue exitoso y mando a iniciar sesion
	app.session.Put(r, "flash", "Registro exitoso, inicie sesion.")
	// Registro y redirecciono a pa pagina de login
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
	// Verifico que las credenciales sean correctas, si no es asi aviso
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
	// tomo el id de usuario para mantener la sesion abierta
	app.session.Put(r, "userID", id)
	// Redireccion a ver turnos.
	http.Redirect(w, r, "/user/turno", http.StatusSeeOther)
}
func (app *application) logoutUser(w http.ResponseWriter, r *http.Request) {
	// Elimino el id de sesion para poder cerrarla
	app.session.Remove(r, "userID")
	// Aviso que se cerro la sesion y redirecciono a la pagina raiz
	app.session.Put(r, "flash", "Ha cerrado sesion con exito")
	http.Redirect(w, r, "/", 303)
}

func (app *application) turnoList(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "turn.page.tmpl", &templateData{
		Form: forms.New(nil),
	})
}

func (app *application) turnoSave(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	form := forms.New(r.PostForm)
	form.Required("fecha")
	err = app.turnos.Insert(form.Get("slider_example_2"))

	app.render(w, r, "turn.page.tmpl", &templateData{
		Form: forms.New(nil),
	})
}

func ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}

//func onSignIn(googleUser) {
//const googleJWT = googleUser.getAuthResponse().id_token
//}

func (app *application) userList(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var users []models.User
	result, err := app.users.DB.Query("SELECT * FROM kinetur.users ")
	if err != nil {
		app.serverError(w, err)
	}
	defer result.Close()
	for result.Next() {
		var user models.User
		err := result.Scan(&user.ID, &user.Tipo, &user.Nombre, &user.Apellido, &user.DNI, &user.Domicilio, &user.Email, &user.Password, &user.Created)
		if err != nil {
			app.serverError(w, err)
		}
		users = append(users, user)
	}
	json.NewEncoder(w).Encode(users)
}

func (app *application) userCreate(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		app.serverError(w, err)
	}
	stmt, err := app.users.DB.Prepare("INSERT INTO kinetur.users (tipo,nombre, apellido, dni, domicilio, email, Password, created) VALUES(?,?,?,?,?,?,?, UTC_TIMESTAMP())")
	if err != nil {
		app.serverError(w, err)
	}

	keyVal := make(map[string]string)
	json.Unmarshal(body, &keyVal)
	tipo := keyVal["tipo"]
	nombre := keyVal["nombre"]
	apellido := keyVal["apellido"]
	dni := keyVal["dni"]
	domicilio := keyVal["domicilio"]
	email := keyVal["email"]
	password := keyVal["password"]

	_, err = stmt.Exec(tipo, nombre, apellido, dni, domicilio, email, password)
	if err != nil {
		app.serverError(w, err)
	}
	fmt.Fprintf(w, "Nuevo user creado")
}



func (app *application) userDelete(w http.ResponseWriter, r *http.Request) {
	userId, err := strconv.Atoi(r.URL.Query().Get(":id"))

	stmt, err := app.users.DB.Prepare("DELETE FROM kinetur.users WHERE id = ?")
	if err != nil {
		app.serverError(w, err)
	}
	_, err = stmt.Exec(userId)
	if err != nil {
		app.serverError(w, err)
	}
	fmt.Fprintf(w, "Post with ID = %s was deleted", userId)
}
