package mysql

import (
	"database/sql"
	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
	"regexp"
	"strings"
	"tp-ISA-go.org/kinetur/pkg/models"
)

// PacientesModel tomo la conexion a db
type PacientesModel struct {
	DB *sql.DB
}

// Insert inserto nuevo paciente en la db
func (m *PacientesModel) Insert(dni, nombres, apellidos, direccion, email, password string) error {
	// Create a bcrypt hash of the plain-text password. nahuel1234
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), 6)
	if err != nil {
		return err
	}
	stmt := "INSERT INTO Pacientes (dni,nombres, apellidos,direccion, email, Password) VALUES(?,?,?,?,?,?)"

	_, err = m.DB.Exec(stmt, dni, nombres, apellidos, direccion, email, string(hashPassword))
	if err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok {
			if mysqlErr.Number == 1062 && strings.Contains(mysqlErr.Message, "users_uc_email") {
				return models.ErrDuplicateEmail
			}
		}
	}
	return err
}

// Authenticate Usaremos el método Authenticate para verificar si un paciente existe con
// la dirección de correo electrónico y la contraseña proporcionadas.
//Esto devolverá el ID (DNI) de usuario si lo hacen.
func (m *PacientesModel) Authenticate(email, password string) (int, error) {
	var id int
	var hashedPassword []byte
	row := m.DB.QueryRow("select dni, password from Pacientes where email = email")
	err := row.Scan(&id, &hashedPassword)
	if err == sql.ErrNoRows {
		return 0, models.ErrInvalidCredentials
	} else if err != nil {
		return 0, err
	}
	// Valida el hash del pass, si no coinciden tira errinvalidcredentials
	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return 0, models.ErrInvalidCredentials
	} else if err != nil {
		return 0, err
	}
	// Devuelve el user ID si esta bien autenticado.
	return id, nil
}

//esta funcion la usa el middleware
func (m *PacientesModel) Get(id int) (*models.Pacientes, error) {
	s := &models.Pacientes{}
	stmt := `SELECT dni, nombres, email, created FROM Pacientes WHERE dni = ?`
	err := m.DB.QueryRow(stmt, id).Scan(&s.DNI, &s.Nombres, &s.Email)
	if err == sql.ErrNoRows {
		return nil, models.ErrNoRecord
	} else if err != nil {
		return nil, err
	}
	return s, nil
}
func validEmail(email string) bool {
	// Validar el formato de la direccion de email.
	var re = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	return re.MatchString(email) && len(email) <= 100
}

func validPassword(pass string) bool {
	return len(pass) >= 6
}
func validDNI(dni string) bool {
	return len(dni) == 8
}
