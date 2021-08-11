package main

import (
	"github.com/bmizerany/pat"
	"github.com/justinas/alice"
	"net/http"
)

func (app *application) routes() http.Handler {
	//creacion del middleware que registra todos los "movimientos"
	standardMiddleware := alice.New(app.logRequest)
	//agrego un middleware dinamico para que tome la session. en el otro me quedan los archivos estaticos
	//no surf evita el crsf
	dynamicMiddleware := alice.New(app.session.Enable, app.authenticate)
	//pat sigue el orden
	//es mas complicado al principio pero cuando crezca al applicacion es mas facil manejar asi los logs de errores
	mux := pat.New()
	mux.Get("/", dynamicMiddleware.ThenFunc(app.home))

	mux.Get("/user/signup", dynamicMiddleware.ThenFunc(app.signupUserForm))
	mux.Post("/user/signup", dynamicMiddleware.ThenFunc(app.signupUser))
	mux.Get("/user/login", dynamicMiddleware.ThenFunc(app.loginUserForm))
	mux.Post("/user/login", dynamicMiddleware.ThenFunc(app.loginUser))
	mux.Get("/user/logout", dynamicMiddleware.ThenFunc(app.logoutUser))

	mux.Get("/user/turno", dynamicMiddleware.ThenFunc(app.turnoList))
	//mux.Post("/user/turno", dynamicMiddleware.ThenFunc(app.guardarTurno))

	//routes de la API
	mux.Get("/pacientes", dynamicMiddleware.ThenFunc(app.userList))
	mux.Get("/profesionales", dynamicMiddleware.ThenFunc(app.profesionalesLista))
	mux.Post("/profesionales", http.HandlerFunc(app.addProfesional))
	mux.Del("/profesionales/:id", http.HandlerFunc(app.removeProfesional))
	mux.Get("/especialidades", dynamicMiddleware.ThenFunc(app.especialidadesLista))
	mux.Post("/especialidades", http.HandlerFunc(app.addEspecialidad))
	mux.Del("/especialidades/:id", http.HandlerFunc(app.removeEspecialidad))

	//crea un servidor de archivos estaticos q estan alojados en ./iu/static
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Get("/static/", http.StripPrefix("/static", fileServer))
	//ruta para test de handlre si la saco el test demuetra q falla
	mux.Get("/ping", http.HandlerFunc(ping))

	return standardMiddleware.Then(mux)
}
