package sqlclient

import (
	"database/sql"
	"errors"
	_ "github.com/go-sql-driver/mysql"
	"os"
)

var (
	dbClient SqlClient
)

const (
	goEnvironment = "GO_ENVIRONMENT"
	production = "production"
)

type client struct {
	db *sql.DB
}


type SqlClient interface {
	Query(query string, args ...interface{}) (rows, error)
}

func StartMockServer()  {
	isMocked = true
}

func StopMockServer()  {
	isMocked = false
}

func isProduction() bool {
	return os.Getenv(goEnvironment) == production
}

func Open(driverName, dataSourceName string) (SqlClient, error) {
	if !isProduction() && isMocked{
		dbClient = &clientMock{}
		return dbClient, nil
	}
	if driverName == "" {
		return nil, errors.New("invalid driver name")
	}
	database, err := sql.Open(driverName, dataSourceName)
	if err != nil{
		return nil, err
	}
	dbClient = &client{
		db: database,
	}
	return dbClient, nil
}

func (c *client) Query(query string, args ...interface{}) (rows, error) {
	returnedRows, err := c.db.Query(query, args...)
	if err != nil{
		return nil, err
	}
	result := sqlRows{
		rows: returnedRows,
	}
	return &result, nil
}