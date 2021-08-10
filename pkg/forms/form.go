//Funciones de validacion de formularios

package forms

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"unicode/utf8"
)

// EmailRX expresion regular para el campo de email
var EmailRX = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9]")

// Form creo una estructura de formulario para guardar valores y errores
type Form struct {
	url.Values
	Errors errors
}

// New Inicializo un nuevo formulario
func New(data url.Values) *Form {
	return &Form{
		data,
		errors(map[string][]string{}),
	}
}

// Required Obliga a completar los valores que sean marcados como Required
func (f *Form) Required(fields ...string) {
	for _, field := range fields {
		value := f.Get(field)
		if strings.TrimSpace(value) == "" {
			f.Errors.Add(field, "El campo no puede estar en blanco")
		}
	}
}

// MaxLength Funcion para comprobar el largo maximo del campo del formulario
func (f *Form) MaxLength(field string, d int) {
	value := f.Get(field)
	if value == "" {
		return
	}
	if utf8.RuneCountInString(value) > d {
		f.Errors.Add(field, fmt.Sprintf("El valor ingresado es muy largo (el maximo es %d,)", d))
	}
}

// PermittedValues comprueba que el texto ingresado coincida con valores especificos que estan permitidos
func (f *Form) PermittedValues(field string, opts ...string) {
	value := f.Get(field)
	if value == "" {
		return
	}
	for _, opt := range opts {
		if value == opt {
			return
		}
	}
	f.Errors.Add(field, "El valor ingresado es invalido")
}

// Valid Validacion que devuelve true si no hay errores
func (f *Form) Valid() bool {
	return len(f.Errors) == 0
}

// MinLength Verifica el valor minimo de caracteres del campo
func (f *Form) MinLength(field string, d int) {
	value := f.Get(field)
	if value == "" {
		return
	}
	if utf8.RuneCountInString(value) < d {
		f.Errors.Add(field, fmt.Sprintf("El valor ingresado es muy corto (el minimo es %d)", d))
	}
}

// MatchesPattern comprueba que el texto ingresado coincida con valores especificos de expresiones regulares
func (f *Form) MatchesPattern(field string, pattern *regexp.Regexp) {
	value := f.Get(field)
	if value == "" {
		return
	}
	if !pattern.MatchString(value) {
		f.Errors.Add(field, "This field is invalid")
	}
}
