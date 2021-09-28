package mysql

import (
	"database/sql"
	"fmt"
	"regexp"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
)

// MySQLBackend is a basic backend handling mysql related stuff.
type MySQLBackend struct {
	Username string
	Password string
	Hostname string
	Port     int

	db *sql.DB
}

// Connects to the database
func (m *MySQLBackend) connect() error {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/", m.Username, m.Password, m.Hostname, m.Port))

	// if there is an error opening the connection, handle it
	if err != nil {
		return err
	}

	m.db = db

	return nil
}

// execute runs a single SQL query and doesn't care about its result unless it's an error.
func (m *MySQLBackend) execute(query string, args ...interface{}) error {
	_, err := m.db.Query(query, args...)
	if err != nil {
		return errors.Wrap(err, "SQL query: "+query)
	}

	return nil
}

// Close closes connection to the database
func (m *MySQLBackend) close() error {
	return m.db.Close()
}

// testValue tests string input for unwanted characters
func (m *MySQLBackend) testValue(value string) error {
	matched, err := regexp.MatchString(`^[a-zA-Z0-9_\.]*$`, value)
	if err != nil {
		return errors.Wrap(err, "regexp error")
	}
	if matched {
		return nil
	}

	return errors.New("invalid value")
}

// This is basic escape function used only for passwords
func (m *MySQLBackend) escape(value string) string {
	value = strings.Replace(value, "'", `\'`, -1)
	value = strings.Replace(value, "\"", `\"`, -1)
	return value
}

func (m *MySQLBackend) CreateROUser(user, password, database string) error {
	if m.testValue(user) != nil {
		return errors.New("invalid format of username")
	}
	if m.testValue(database) != nil {
		return errors.New("invalid format of database")
	}

	if err := m.connect(); err != nil {
		return err
	}
	defer m.close()

	sqls := []string{
		"CREATE USER '" + user + "'@'%' IDENTIFIED BY '" + m.escape(password) + "';",
		"GRANT SELECT PRIVILEGES ON " + database + ".* TO '" + user + "'@'%';",
	}

	for _, sql := range sqls {
		err := m.execute(sql)
		if err != nil {
			return err
		}
	}

	return nil
}

func (m *MySQLBackend) CreateUser(user, password, database string) error {
	if m.testValue(user) != nil {
		return errors.New("invalid format of username")
	}
	if m.testValue(database) != nil {
		return errors.New("invalid format of database")
	}

	if err := m.connect(); err != nil {
		return err
	}
	defer m.close()

	sql := "CREATE USER '" + user + "'@'%' IDENTIFIED BY '" + m.escape(password) + "';"
	return m.execute(sql)
}

func (m *MySQLBackend) CreateDatabase(database, owner string, extensions []string) error {
	if m.testValue(owner) != nil {
		return errors.New("invalid format of owner")
	}
	if m.testValue(database) != nil {
		return errors.New("invalid format of database")
	}

	if err := m.connect(); err != nil {
		return err
	}
	defer m.close()

	sql := "CREATE DATABASE " + database + ";"
	err := m.execute(sql)
	if err != nil {
		return err
	}

	sql = "GRANT ALL PRIVILEGES ON " + database + ".* TO '" + owner + "'@'%';"
	err = m.execute(sql)
	if err != nil {
		return err
	}

	sql = "FLUSH PRIVILEGES;"

	return m.execute(sql)
}

func (m *MySQLBackend) ChangePassword(user, password string) error {
	if m.testValue(user) != nil {
		return errors.New("invalid format of user")
	}

	if err := m.connect(); err != nil {
		return err
	}
	defer m.close()

	sql := "SET PASSWORD FOR '" + user + "'@'%' = PASSWORD('" + m.escape(password) + "');"
	return m.execute(sql)
}

func (m *MySQLBackend) DropUser(user string) error {
	if m.testValue(user) != nil {
		return errors.New("invalid format of user")
	}

	if err := m.connect(); err != nil {
		return err
	}
	defer m.close()

	sql := "DROP USER '" + user + "';"
	return m.execute(sql)
}

func (m *MySQLBackend) DropDatabase(database string) error {
	if m.testValue(database) != nil {
		return errors.New("invalid format of database")
	}

	if err := m.connect(); err != nil {
		return err
	}
	defer m.close()

	sql := "DROP DATABASE " + database + ";"
	return m.execute(sql)
}
