package main

import (
	"context"
	"fmt"
	_ "fmt"
	"github.com/justinas/nosurf"
	"net/http"
	"tp-ISA-go.org/kinetur/pkg/models"
)

// Le agrega a cada respuesta que emite el servidor parametros de aseguramiento de header.
func secureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//Este codigo se ejecuta antes de llegar al Application Handler!!!
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("X-Frame-Options", "deny")
		next.ServeHTTP(w, r)
		//El codigo aca se ejecuta despues de pasar por el Application handler
	})
}

//funcion que muestra el log por terminal
func (app *application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.infoLog.Printf("%s - %s %s %s", r.RemoteAddr, r.Proto, r.Method, r.URL)
		next.ServeHTTP(w, r)
	})
}

func (app *application) requireAuthenticatedUser(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Si el usuario no esta autenticado se lo redirecciona al login
		// desde el middleware, sin llegar a ejecutar el handler
		if app.authenticatedUser(r) == nil {
			http.Redirect(w, r, "/user/login", 302)
			return
		}
		// si esta logueado pasa al handler
		next.ServeHTTP(w, r)
	}
}

//cuando falla la aplicacion muestra internal server error
func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
				app.serverError(w, fmt.Errorf(" %s", err))
			}
		}()
		next.ServeHTTP(w, r)
	})
}

//no-surf previene CSRF (cross-site request forgery)
func noSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)
	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   true,
	})
	return csrfHandler
}

//chequea que el ID de usuario exista en la sesion
func (app *application) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		exists := app.session.Exists(r, "userID")
		if !exists {
			next.ServeHTTP(w, r)
			return
		}
		user, err := app.pacientes.Get(app.session.GetInt(r, "userID"))
		if err == models.ErrNoRecord {
			app.session.Remove(r, "userID")
			next.ServeHTTP(w, r)
			return
		} else if err != nil {
			app.serverError(w, err)
			return
		}
		ctx := context.WithValue(r.Context(), contextKeyUser, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (app *application) withCORS(next http.Handler) http.Handler {
	// para mostrar por la salida de log cada request que se le haga al server
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// La peticion debe venir de un origen determinado, no vale *
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:4000")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Access-Control-Allow-Headers, Authorization, X-Requested-With")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		// Stop here for a Preflighted OPTIONS request.
		next.ServeHTTP(w, r)
	})
}

func (app *application) authenticateUser(r *http.Request) bool {
	var email = app.session.GetString(r, "email")
	if len(email) <= 0 {
		return false
	}
	return true
}

func (app *application) restrictedEndpoint(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !app.authenticateUser(r) {
			app.clientError(w, 401) //unauthorized
			return
		}
		next.ServeHTTP(w, r)
	})

}
