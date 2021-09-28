package pgsql

import (
	"database/sql"
	"fmt"
	"regexp"
	"strings"

	_ "github.com/lib/pq"
	"github.com/pkg/errors"
)

// PGSQLBackend is a basic backend handling pgsql related stuff.
type PGSQLBackend struct {
	Username string
	Password string
	Hostname string
	Port     int

	db *sql.DB
}

// Connects to the database
// Database with same name as the username has to exist
func (p *PGSQLBackend) connect(database string) error {
	db, err := sql.Open("postgres", fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", p.Hostname, p.Port, p.Username, p.Password, database))

	// if there is an error opening the connection, handle it
	if err != nil {
		return err
	}

	p.db = db

	return nil
}

// execute runs a single SQL query and doesn't care about its result unless it's an error.
func (p *PGSQLBackend) execute(query string, args ...interface{}) error {
	_, err := p.db.Query(query, args...)
	if err != nil {
		return errors.Wrap(err, "SQL query: "+query)
	}

	return nil
}

// testValue tests string input for unwanted characters
func (p *PGSQLBackend) testValue(value string) error {
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
func (p *PGSQLBackend) escape(value string) string {
	value = strings.Replace(value, "'", `\'`, -1)
	value = strings.Replace(value, "\"", `\"`, -1)
	return value
}

// Close closes connection to the database
func (p *PGSQLBackend) close() error {
	return p.db.Close()
}

func (p *PGSQLBackend) CreateUser(user, password, database string) error {
	if p.testValue(user) != nil {
		return errors.New("invalid format of username")
	}
	if p.testValue(database) != nil {
		return errors.New("invalid format of database")
	}

	if err := p.connect(p.Username); err != nil {
		return err
	}
	defer p.close()

	sql := "CREATE USER " + p.escape(user) + " WITH PASSWORD '" + password + "';"
	return p.execute(sql)
}

func (p *PGSQLBackend) CreateROUser(user, password, database string) error {
	sqls := []string{
		fmt.Sprintf("CREATE USER %s WITH PASSWORD '%s';", user, password),
		fmt.Sprintf("GRANT CONNECT ON DATABASE %s TO %s;", database, user),
		fmt.Sprintf("GRANT USAGE ON SCHEMA %s TO %s;", database, user),
		fmt.Sprintf("GRANT SELECT ON ALL TABLES IN SCHEMA %s TO %s;", database, user),
		fmt.Sprintf("ALTER DEFAULT PRIVILEGES IN SCHEMA %s GRANT SELECT ON TABLES TO %s;", database, user),
	}

	for _, sql := range sqls {
		err := p.execute(sql)
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *PGSQLBackend) CreateDatabase(database, owner string, extensions []string) error {
	if p.testValue(owner) != nil {
		return errors.New("invalid format of owner")
	}
	if p.testValue(database) != nil {
		return errors.New("invalid format of database")
	}
	for _, extension := range extensions {
		if p.testValue(extension) != nil {
			return errors.New("invalid format of extension")
		}
	}

	if err := p.connect(p.Username); err != nil {
		return err
	}

	sql := "CREATE DATABASE " + database + " OWNER " + owner + ";"
	err := p.execute(sql)
	if err != nil {
		return err
	}
	p.close()

	if err := p.connect(database); err != nil {
		return err
	}
	defer p.close()

	sql = "CREATE SCHEMA " + database + ";"
	err = p.execute(sql)
	if err != nil {
		return err
	}

	sql = "ALTER SCHEMA " + database + " OWNER TO " + owner + ";"
	err = p.execute(sql)
	if err != nil {
		return err
	}

	for _, extension := range extensions {
		sql := "CREATE EXTENSION " + extension + " SCHEMA " + database + ";"
		err := p.execute(sql)
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *PGSQLBackend) ChangePassword(user, password string) error {
	if p.testValue(user) != nil {
		return errors.New("invalid format of username")
	}

	if err := p.connect(p.Username); err != nil {
		return err
	}
	defer p.close()

	sql := "ALTER USER " + p.escape(user) + " PASSWORD '" + password + "';"
	return p.execute(sql)
}

func (p *PGSQLBackend) DropUser(user string) error {
	if p.testValue(user) != nil {
		return errors.New("invalid format of username")
	}

	if err := p.connect(p.Username); err != nil {
		return err
	}
	defer p.close()

	sql := "DROP OWNED BY " + p.escape(user) + " CASCADE;"
	err := p.execute(sql)
	if err != nil {
		return err
	}

	sql = "DROP ROLE " + p.escape(user) + ";"
	err = p.execute(sql)
	return err
}

func (p *PGSQLBackend) DropDatabase(database string) error {
	// Is this needed?
	// self.get_connection().set_isolation_level(psycopg2.extensions.ISOLATION_LEVEL_AUTOCOMMIT)

	if p.testValue(database) != nil {
		return errors.New("invalid format of database")
	}

	if err := p.connect(p.Username); err != nil {
		return err
	}
	defer p.close()

	sql := "SELECT pg_terminate_backend(pid) FROM pg_stat_activity WHERE datname = '" + database + "';"
	err := p.execute(sql)
	if err != nil {
		return err
	}
	sql = "DROP DATABASE " + database + ";"
	err = p.execute(sql)
	return err
}
