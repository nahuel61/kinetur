package mysql

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"regexp"
	"tp-ISA-go.org/kinetur/pkg/models"
)

type ProfesionalesModel struct {
	DB *sql.DB
}

func (m *ProfesionalesModel) Insert(dni, nombres, apellidos, especialidad string) error {
	stmt := "INSERT INTO kinetur.Profesional (dni,nombres, apellidos, especialidad_id) VALUES(?,?,?,?)"

	_, err := m.DB.Exec(stmt, dni, nombres, apellidos, especialidad)
	if err != nil {
		return err
	}

	return err
}
func (m *ProfesionalesModel) Get(id int) (*models.Profesionales, error) {
	s := &models.Profesionales{}
	stmt := `SELECT * FROM kinetur.Profesional`
	err := m.DB.QueryRow(stmt, id).Scan(&s.Id, &s.Nombres, &s.Apellidos, &s.Especialidad)
	if err == sql.ErrNoRows {
		return nil, models.ErrNoRecord
	} else if err != nil {
		return nil, err
	}
	return s, nil
}

func (m *ProfesionalesModel) Del(id int) error {
	stmt := "DELETE from kinetur.Profesional where id=?"

	_, err := m.DB.Exec(stmt, id)
	if err != nil {
		return err
	}
	return err
}

func validUsername(username string) bool {
	// Validar el nombre de usuario de acuerdo a ese regex.
	var re = regexp.MustCompile("^[[:alpha:] '-]+$")
	return re.MatchString(username) && len(username) <= 100 && len(username) >= 6
}
