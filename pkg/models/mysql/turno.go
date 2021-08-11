package mysql

import (
	"database/sql"
)

// TurnoModel tomo la conexion a db
type TurnoModel struct {
	DB *sql.DB
}

// inserto turno en la db
func (m *TurnoModel) Insert(paciente, profesional, fecha string) error {

	stmt := "INSERT INTO kinetur.Citas (paciente_DNI, profesional_id, fecha) VALUES(?,?,?)"
	_, err := m.DB.Exec(stmt, paciente, profesional, fecha)
	if err != nil {
		panic(err)
	}
	return nil
}
