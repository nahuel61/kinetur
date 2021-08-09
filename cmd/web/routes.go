package main

import (
	"github.com/bmizerany/pat"
	"github.com/justinas/alice"
	"net/http"
)

func (app *application) routes() http.Handler {
	//creacion del middleware que registra todos los "movimientos"
	standardMiddleware := alice.New(app.logRequest, secureHeaders)
	//agrego un middleware dinamico para que tome la session. en el otro me quedan los archivos estaticos
	dynamicMiddleware := alice.New(app.session.Enable, noSurf, app.authenticate)
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

	//routes de la API
	mux.Get("/pacientes", dynamicMiddleware.ThenFunc(app.userList))
	mux.Post("/user", dynamicMiddleware.ThenFunc(app.userCreate))
	mux.Del("/user/{dni}", dynamicMiddleware.ThenFunc(app.userDelete))
	mux.Get("/profesionales", dynamicMiddleware.ThenFunc(app.profesionalesLista))

	//crea un servidor de archivos estaticos q estan alojados en ./iu/static
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Get("/static/", http.StripPrefix("/static", fileServer))

	return standardMiddleware.Then(mux)
}
