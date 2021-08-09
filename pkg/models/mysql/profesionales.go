package mysql

import "database/sql"

type ProfesionalesModel struct {
	DB *sql.DB
}

/*
func (m *ProfesionalesModel) Get(id int) (*models.Profesional, error) {
	s := &models.Profesional{}
	stmt := `SELECT * FROM kinetur.Profesional`
	err := m.DB.QueryRow(stmt, id).Scan(&s.Id, &s.Nombres, &s.Apellidos,&s.Especialidad)
	if err == sql.ErrNoRows {
		return nil, models.ErrNoRecord
	} else if err != nil {
		return nil, err
	}
	return s, nil
}*/
