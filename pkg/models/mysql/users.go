package mysql

import (
	"database/sql"
	"tp-ISA-go.org/kinetur/pkg/models"
)

//tomo la conexion a db
type UsersModel struct {
	DB *sql.DB
}

// inserto nuevo usuario en la db
func (m *UsersModel) Insert(nombre, apellido, dni string) (int, error) {
	stmt := "INSERT INTO users (nombre, apellido,dni) VALUES (? , ? ,?)"

	result, err := m.DB.Exec(stmt, nombre, apellido, dni)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

//devuelve un usuario basado en su id
func (m *UsersModel) GET(id int) (*models.Users, error) {
	stmt := "SELECT id, nombre, apellido, dni FROM users"

	row := m.DB.QueryRow(stmt) //devuelve un puntero a sql.row que contiene el resultado de la consulta a la db
	s := &models.Users{}       //inicializo un puntero a una structura vacia

	err := row.Scan(&s.ID, &s.Nombre, &s.Apellido, &s.DNI)
	if err == sql.ErrNoRows {
		return nil, models.ErrNoRecord
	} else if err != nil {
		return nil, err
	}
	return s, nil

}

//devuelve los 10 ultimos usuarios
func (m *UsersModel) Latest() ([]*models.Users, error) {
	return nil, nil
}
