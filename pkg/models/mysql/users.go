package mysql

import (
	"database/sql"
	"strings"
	"tp-ISA-go.org/kinetur/pkg/models"

	"github.com/go-sql-driver/mysql" // New import
	"golang.org/x/crypto/bcrypt"     // New import
)

//tomo la conexion a db
type UserModel struct {
	DB *sql.DB
}

// inserto nuevo usuario en la db
func (m *UserModel) Insert(tipo, nombre, apellido, dni, domicilio, email, password string) error {
	// Create a bcrypt hash of the plain-text password.
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), 8)
	if err != nil {
		return err
	}
	stmt := "INSERT INTO users (tipo,nombre, apellido, dni, domicilio, email, Password, created) VALUES(?,?,?,?,?,?,?, UTC_TIMESTAMP())"
	_, err = m.DB.Exec(stmt, tipo, nombre, apellido, dni, domicilio, email, string(hashPassword))
	if err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok {
			if mysqlErr.Number == 1062 && strings.Contains(mysqlErr.Message, "u") {
				return models.ErrDuplicateEmail
			}
		}
	}
	return err
}

// We'll use the Authenticate method to verify whether a user exists with
// the provided email address and password. This will return the relevant
// user ID if they do.
func (m *UserModel) Authenticate(email, password string) (int, error) {
	return 0, nil
}

//devuelve un usuario basado en su id
func (m *UserModel) Get(id int) (*models.User, error) {

	return nil, nil

}
