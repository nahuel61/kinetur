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

func ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
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
	form.Required("dni", "nombres", "apellidos", "direccion", "email", "password")
	form.MatchesPattern("email", forms.EmailRX)
	form.MinLength("dni", 8)
	form.MinLength("password", 6)
	// si hay un error vuelvo a mostrar el formulario
	if !form.Valid() {
		app.render(w, r, "signup.page.tmpl", &templateData{Form: form})
		return
	}
	// creo un nuevo registro en la base de datos, si se repite el email avisa
	err = app.pacientes.Insert(form.Get("dni"), form.Get("nombres"), form.Get("apellidos"), form.Get("direccion"), form.Get("email"), form.Get("password"))
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
	// Registro y redirecciono a la pagina de login
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
	id, err := app.pacientes.Authenticate(form.Get("email"), form.Get("password"))
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

func (app *application) userList(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var users []models.Pacientes
	result, err := app.pacientes.DB.Query("SELECT * FROM kinetur.Pacientes ")
	if err != nil {
		app.serverError(w, err)
	}
	defer result.Close()
	for result.Next() {
		var user models.Pacientes
		err := result.Scan(&user.DNI, &user.Nombres, &user.Apellidos, &user.Direccion, &user.Email, &user.Password)
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
	stmt, err := app.pacientes.DB.Prepare("INSERT INTO kinetur.Pacientes (dni,nombres, apellidos,direccion, email, Password) VALUES(?,?,?,?,?,?)")
	if err != nil {
		app.serverError(w, err)
	}

	keyVal := make(map[string]string)
	json.Unmarshal(body, &keyVal)
	DNI := keyVal["dni"]
	nombres := keyVal["nombres"]
	apellidos := keyVal["apellidos"]
	direccion := keyVal["direccion"]
	email := keyVal["email"]
	password := keyVal["password"]

	_, err = stmt.Exec(DNI, nombres, apellidos, direccion, email, password)
	if err != nil {
		app.serverError(w, err)
	}
	fmt.Fprintf(w, "Nuevo paciente creado")
}

func (app *application) userDelete(w http.ResponseWriter, r *http.Request) {
	userId, err := strconv.Atoi(r.URL.Query().Get(":dni"))

	stmt, err := app.pacientes.DB.Prepare("DELETE FROM kinetur.Pacientes WHERE DNI = ?")
	if err != nil {
		app.serverError(w, err)
	}
	_, err = stmt.Exec(userId)
	if err != nil {
		app.serverError(w, err)
	}
	fmt.Fprintf(w, "Paciente con DNI = %c fue eliminado", userId)
}

//--------------------------FUNCIONES DE LOS PROFESIONALES-------------------//

func (app *application) profesionalesLista(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var profesionales []models.Profesionales
	result, err := app.profesional.DB.Query("SELECT * FROM kinetur.Profesional ")
	if err != nil {
		app.serverError(w, err)
	}
	defer result.Close()
	for result.Next() {
		var prof models.Profesionales
		err := result.Scan(&prof.Id, &prof.DNI, &prof.Nombres, &prof.Apellidos, &prof.Especialidad)
		if err != nil {
			app.serverError(w, err)
		}
		profesionales = append(profesionales, prof)
	}
	json.NewEncoder(w).Encode(profesionales)
}

func (app *application) addProfesional(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		app.serverError(w, err)
		return
	}
	stmt, err := app.profesional.DB.Prepare("INSERT INTO kinetur.Profesional (DNI,nombres, apellidos, especialidad_id) VALUES(?,?,?,?)")
	if err != nil {
		app.serverError(w, err)
	}
	keyVal := make(map[string]string)
	json.Unmarshal(body, &keyVal)
	DNI := keyVal["dni"]
	nombres := keyVal["nombres"]
	apellidos := keyVal["apellidos"]
	especialidadId := keyVal["especialidad"]

	_, err = stmt.Exec(DNI, nombres, apellidos, especialidadId)
	if err != nil {
		app.serverError(w, err)
	}
	fmt.Fprintf(w, "Nuevo profesional agregado")

}

func (app *application) removeProfesional(w http.ResponseWriter, r *http.Request) {
	// Manejador que dada una peticion con el id en la URI, elimina al usuario y devuelve un 200 vacío.
	userID, err := strconv.Atoi(r.URL.Query().Get(":id"))
	_, err = app.profesional.DB.Exec("delete from kinetur.Profesional where id = ?", userID)
	if err != nil && err.Error() == "record not found" {
		app.clientError(w, 404)
		return
	} else if err != nil {
		app.serverError(w, err)
		return
	}
	fmt.Fprintf(w, "Profesional eliminado")
}

//func onSignIn(googleUser) {
//const googleJWT = googleUser.getAuthResponse().id_token
//}
/*func (m *ProfesionalesModel) Insert(DNI int, nombres string, apellidos string, especialidad int) (int,error) {
	// Create a bcrypt hash of the plain-text password. nahuel1234

	stmt := "INSERT INTO kinetur.Profesional (DNI,nombres, apellidos, especialidad_id) VALUES(?,?,?,?)"

	result, err := m.DB.Exec(stmt, DNI, nombres, apellidos, especialidad )
	if err != nil {
		return 0, err
	}
	id, err:= result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil

}
*/
