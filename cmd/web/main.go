package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"
	"os"

	"tp-ISA-go.org/kinetur/pkg/models/mysql"

	_ "github.com/go-sql-driver/mysql"
)

type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
	users    *mysql.UsersModel //con esto permito que este disponible para el handler
}

func main() {
	//defino la direccion default
	dsn := flag.String("dsn", "root:admin@/kinetur?parseTime=true", "Mysql data")
	addr := flag.String("addr", ":4000", "HTTP network address")
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

	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
		users:    &mysql.UsersModel{DB: db},
	}

	//configuro strcutura del http.server
	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	//escribo los mensajes de info y error en cada uno de los logs
	infoLog.Printf("Starting server on %s", *addr)
	err = srv.ListenAndServe()
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
