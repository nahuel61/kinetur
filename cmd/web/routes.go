package main

import (
	"github.com/bmizerany/pat"
	"github.com/justinas/alice"
	"net/http"
)

func (app *application) routes() http.Handler {
	//creacion del middleware que registra todos los "movimientos"
	standardMiddleware := alice.New(app.logRequest, secureHeaders)
	//pat sigue el orden
	//es mas complicado al principio pero cuando crezca al applicacion es mas facil manejar asi los logs de errores
	mux := pat.New()
	mux.Get("/", http.HandlerFunc(app.home))

	mux.Get("/user/signup", http.HandlerFunc(app.signupUserForm))
	mux.Post("/user/signup", http.HandlerFunc(app.signupUser))
	mux.Get("/user/login", http.HandlerFunc(app.loginUserForm))
	mux.Post("/user/login", http.HandlerFunc(app.loginUser))
	mux.Post("/user/logout", http.HandlerFunc(app.logoutUser))

	//crea un servidor de archivos estaticos q estan alojados en ./iu/static
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Get("/static/", http.StripPrefix("/static", fileServer))

	return standardMiddleware.Then(mux)
}
