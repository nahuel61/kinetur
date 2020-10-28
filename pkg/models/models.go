package models

import (
	"errors"
	"time"
)

var (
	ErrNoRecord           = errors.New("models: no matching record found")
	ErrInvalidCredentials = errors.New("models: invalid credentials")
	ErrDuplicateEmail     = errors.New("models: duplicate email")
)

type User struct {
	ID        int       `json:"id"`
	Tipo      string    `json:"tipo"`
	Nombre    string    `json:"nombre"`
	Apellido  string    `json:"apellido"`
	DNI       int       `json:"dni"`
	Domicilio string    `json:"domicilio"`
	Email     string    `json:"email"`
	Password  []byte    `json:"password"`
	Created   time.Time `json:"created"`
}

type Turno struct {
	ID    int
	Fecha time.Time
}
