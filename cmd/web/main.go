package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

func main() {
	//defino la direccion default
	addr := flag.String("addr", ":4000", "HTTP network address")
	//con flag.Parse() hago que tome la direccion por default cuando corro la aplicacion
	//go run cmd/web/ -addr=":xxx"
	flag.Parse()

	// voy a crear log diferenciados, de eventos y de errores
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)


	mux := http.NewServeMux()
	mux.HandleFunc("/", home)
	mux.HandleFunc("/register", crearUsuario)
	mux.HandleFunc("/login", iniciarSesion)


	//crea un servidor de archivos estaticos q estan alojados en ./iu/static
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))



	//escribo los mensajes de info y error en cada uno de los logs
	infoLog.Printf("Starting server on %s", *addr)
	err := http.ListenAndServe(*addr, mux)
	errorLog.Fatal(err)
}
