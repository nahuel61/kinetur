package mysql

import (
	"database/sql"
	"tp-ISA-go.org/kinetur/pkg/models"
)

type EspecialidadModel struct {
	DB *sql.DB
}

func (m *EspecialidadModel) Insert(nombre string) error {
	stmt := "INSERT INTO kinetur.Especialidades (nombre) VALUES(?)"

	_, err := m.DB.Exec(stmt, nombre)
	if err != nil {
		return err
	}
	return err
}
func (m *EspecialidadModel) Get(id int) (*models.Especialidades, error) {
	s := &models.Especialidades{}
	stmt := `SELECT * FROM kinetur.Especialidades`
	err := m.DB.QueryRow(stmt, id).Scan(&s.Id, &s.Nombre)
	if err == sql.ErrNoRows {
		return nil, models.ErrNoRecord
	} else if err != nil {
		return nil, err
	}
	return s, nil
}
