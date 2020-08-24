package main

import (
	"log"
	"net/http"
)

func home(w http.ResponseWriter, r *http.Request) {
	//comprueba si concide exactamente si no devuelve 404
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	w.Write([]byte("Hola, vamos a sacar un turno para tu rehabilitacion!!!"))
}

func registrarse(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hola, vamos registrate"))
}
func identificarse(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hola, vamos iniciar sesion"))
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", home)
	mux.HandleFunc("/newuser", registrarse)
	mux.HandleFunc("/login", identificarse)
	log.Println("starting server on :4000")
	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}
