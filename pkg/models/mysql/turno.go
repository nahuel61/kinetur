package mysql

import (
	"database/sql"
)

// TurnoModel tomo la conexion a db
type TurnoModel struct {
	DB *sql.DB
}

/* inserto turno en la db
func (m *TurnoModel) Insert(fecha string) error {

	stmt := "INSERT INTO turno (fecha) VALUES(?)"
	err, _ := m.DB.Exec(stmt, fecha)
	if err != nil{
		panic(err)
	}
	return nil
}

func (m *TurnoModel) Get(id int) (*models.Turno, error) {
	s := &models.Turno{}
	stmt := `SELECT fecha FROM turno WHERE id = ?`
	err := m.DB.QueryRow(stmt, id).Scan(&s.ID, &s.Fecha)
	if err == sql.ErrNoRows {
		return nil, models.ErrNoRecord
	} else if err != nil {
		return nil, err
	}
	return s, nil
}

*/
