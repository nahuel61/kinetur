package main

import (
	"crypto/tls"
	"database/sql"
	"flag"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golangcollege/sessions"
	//_ "golang.org/x/oauth2/google"
	//_ "google.golang.org/api/calendar/v3"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"
	"tp-ISA-go.org/kinetur/pkg/models/mysql"
)

type contextKey string

var contextKeyUser = contextKey("paciente")

type application struct {
	errorLog      *log.Logger
	infoLog       *log.Logger
	pacientes     *mysql.PacientesModel //con esto permito que este disponible para el handler
	profesional   *mysql.ProfesionalesModel
	templateCache map[string]*template.Template
	session       *sessions.Session //agrego sesion a la struc
	turnos        *mysql.TurnoModel
}

func main() {

	//defino la direccion default
	dsn := flag.String("dsn", "root:admin@/kinetur?parseTime=true", "Mysql data")
	addr := flag.String("addr", ":4000", "HTTPS network address")

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
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			return
		}
	}(db)

	// inicializo el cache de templates
	templateCache, err := newTemplateCache("./ui/html/")
	if err != nil {
		errorLog.Fatal(err)
	}
	//creo un nuevo session manager y le paso la clave secret como parametro. la sesion expira a las 6 horas
	session := sessions.New([]byte(*secret))
	session.Lifetime = 10 * time.Minute
	session.Secure = true
	session.Persist = false //cierra el explorador y se pierda la cookie
	session.SameSite = http.SameSiteStrictMode

	app := &application{
		errorLog:      errorLog,
		infoLog:       infoLog,
		session:       session,
		pacientes:     &mysql.PacientesModel{DB: db},
		profesional:   &mysql.ProfesionalesModel{DB: db},
		templateCache: templateCache,
		turnos:        &mysql.TurnoModel{DB: db},
	}

	// Initialize a tls.Config struct to hold the non-default TLS settings we want the server to use.
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
		ReadTimeout:  5 * time.Second, //previene ataques de cliente lento (mantiene abierta la conexxcion la mas que pueda)
		WriteTimeout: 10 * time.Second,
	}

	//escribo los mensajes de info y error en cada uno de los logs
	infoLog.Printf("Inicio el servidor en https://localhost%s", *addr)
	//agrego certificacion de seguridad tls para iniciar servidor https
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
