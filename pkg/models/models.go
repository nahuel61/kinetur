package models

import (
	"errors"
	"time"
	_ "time"
)

var (
	ErrNoRecord           = errors.New("models: no se encontró ningún registro coincidente")
	ErrInvalidCredentials = errors.New("models: credenciales invalidas")
	ErrDuplicateEmail     = errors.New("models: email duplicado")
)

type Pacientes struct {
	DNI       string `json:"dni"`
	Nombres   string `json:"nombre"`
	Apellidos string `json:"apellido"`
	Direccion string `json:"direccion"`
	Email     string `json:"email"`
	Password  []byte `json:"password"`
}

type Especialidades struct {
	Id      int    `json:"id"`
	Nombres string `json:"nombre"`
}

type Profesional struct {
	Id           int    `json:"id"`
	DNI          string `json:"dni"`
	Nombres      string `json:"nombre"`
	Apellidos    string `json:"apellido"`
	Direccion    string `json:"direccion"`
	Especialidad string `json:"especialidad"`
}

type Horarios struct {
	Id         int    `json:"id"`
	HoraInicio string `json:"hora_inicio"`
	HoraFin    string `json:"hora_fin"`
	Cupos      string `json:"cupos"`
}

type Dias struct {
	Dia           string `json:"dia"`
	ProfesionalId int    `json:"profesional_id"`
	HorarioId     string `json:"horario_id"`
}

type Citas struct {
	Id            int       `json:"id"`
	PacienteDNI   string    `json:"paciente_dni"`
	ProfesionalId string    `json:"profesional_id"`
	Fecha         time.Time `json:"fecha"`
}
