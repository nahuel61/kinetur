package main

import (
	"net/http"
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

func (app *application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.infoLog.Printf("%s - %s %s %s", r.RemoteAddr, r.Proto, r.Method, r.URL)
		next.ServeHTTP(w, r)
	})
}
