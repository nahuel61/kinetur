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
	ID        int
	Tipo      string
	Nombre    string
	Apellido  string
	DNI       int
	Domicilio string
	Email     string
	Password  []byte
	Created   time.Time
}
