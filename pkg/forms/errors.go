package forms

//Funciones para definir tipos de errores de las validaciones de formularios

type errors map[string][]string

// Add Implemento el metodo Add() para agregar un mensaje de rror a un campo dadoa map
func (e errors) Add(field, message string) {
	e[field] = append(e[field], message)
}

// Get Implemento el metodo Get() para recuperar el primer mensaje de errar para un campo dado
func (e errors) Get(field string) string {
	es := e[field]
	if len(es) == 0 {
		return ""
	}
	return es[0]
}
