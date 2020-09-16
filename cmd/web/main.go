package main

import (
	"crypto/tls"
	"database/sql"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golangcollege/sessions" // New import
	"tp-ISA-go.org/kinetur/pkg/models/mysql"

	_ "github.com/go-sql-driver/mysql"
)

type application struct {
	errorLog      *log.Logger
	infoLog       *log.Logger
	users         *mysql.UserModel //con esto permito que este disponible para el handler
	templateCache map[string]*template.Template
	session       *sessions.Session //agrego sesion a la struc
}

func main() {

	//defino la direccion default
	dsn := flag.String("dsn", "root:admin@/kinetur?parseTime=true", "Mysql data")
	addr := flag.String("addr", ":4000", "HTTP network address")
	//agrego autenticacion
	secret := flag.String("secret", "s6Ndh+pPbnzHbS*+9Pk8qGWhTzbpa@ge", "Secret key")

	//con flag.Parse() hago que tome la direccion por default cuando corro la aplicacion
	flag.Parse()

	// voy a crear log diferenciados, de eventos y de errores
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()

	// Initialize a new template cache...
	templateCache, err := newTemplateCache("./ui/html/")
	if err != nil {
		errorLog.Fatal(err)
	}
	//creo un nuevo session manager y le paso la clave secret como parametro. la sesion expira a las 12 horas
	session := sessions.New([]byte(*secret))
	session.Lifetime = 12 * time.Hour
	session.Secure = true

	// And add it to the application dependencies.
	app := &application{
		errorLog:      errorLog,
		infoLog:       infoLog,
		users:         &mysql.UserModel{DB: db},
		templateCache: templateCache,
		session:       session,
	}

	// Initialize a tls.Config struct to hold the non-default TLS settings we w
	// the server to use.
	tlsConfig := &tls.Config{
		PreferServerCipherSuites: true,
		CurvePreferences:         []tls.CurveID{tls.X25519, tls.CurveP256},
	}

	//configuro strcutura del http.server
	srv := &http.Server{
		Addr:         *addr,
		ErrorLog:     errorLog,
		Handler:      app.routes(),
		TLSConfig:    tlsConfig,
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	//escribo los mensajes de info y error en cada uno de los logs
	infoLog.Printf("Inicio el servidor en https://localhost%s", *addr)
	//agrego certificacion de seguridad tls
	err = srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
	errorLog.Fatal(err)
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
