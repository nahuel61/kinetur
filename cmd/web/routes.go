package main

import "net/http"

func (app *application) routes() *http.ServeMux {
	//es mas complicado al principio pero cuando crezca al applicacion es mas facil manejar asi los logs de errores
	mux := http.NewServeMux()
	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/register", app.crearUsuario)
	mux.HandleFunc("/login", app.iniciarSesion)
	mux.HandleFunc("/user", app.mostrarUsuario)

	//crea un servidor de archivos estaticos q estan alojados en ./iu/static
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	return mux
}
