package mysql

import (
	"database/sql"
	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
	"strings"
	"tp-ISA-go.org/kinetur/pkg/models"
)

//tomo la conexion a db
type UserModel struct {
	DB *sql.DB
}

// inserto nuevo usuario en la db
func (m *UserModel) Insert(tipo, nombre, apellido, dni, domicilio, email, password string) error {
	// Create a bcrypt hash of the plain-text password. nahuel1234
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), 8)
	if err != nil {
		return err
	}
	stmt := "INSERT INTO users (tipo,nombre, apellido, dni, domicilio, email, Password, created) VALUES(?,?,?,?,?,?,?, UTC_TIMESTAMP())"

	_, err = m.DB.Exec(stmt, tipo, nombre, apellido, dni, domicilio, email, string(hashPassword))
	if err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok {
			if mysqlErr.Number == 1062 && strings.Contains(mysqlErr.Message, "users_uc_email") {
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
	var id int
	var hashedPassword []byte
	row := m.DB.QueryRow("select id, password from users where email = email")
	err := row.Scan(&id, &hashedPassword)
	if err == sql.ErrNoRows {
		return 0, models.ErrInvalidCredentials
	} else if err != nil {
		return 0, err
	}
	// Valida el hash del pass
	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return 0, models.ErrInvalidCredentials
	} else if err != nil {
		return 0, err
	}
	// Devuelve el user ID si esta bien autenticado.
	return id, nil
}

func (m *UserModel) Get(id int) (*models.User, error) {
	s := &models.User{}
	stmt := `SELECT id, nombre, email, created FROM users WHERE id = ?`
	err := m.DB.QueryRow(stmt, id).Scan(&s.ID, &s.Nombre, &s.Email, &s.Created)
	if err == sql.ErrNoRows {
		return nil, models.ErrNoRecord
	} else if err != nil {
		return nil, err
	}
	return s, nil
}
