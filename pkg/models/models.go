package models

import (
	"errors"
	_ "time"
)

var ErrNoRecord = errors.New("models: no matching record found")

type Users struct {
	ID       int
	Nombre   string
	Apellido string
	DNI      int
}
