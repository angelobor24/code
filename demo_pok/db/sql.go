package db

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

// initialize the DB structure, specifying the db drive, name and db configuration
func InitializeDB(dbType string, dbName string, enabledFlag string) (*sql.DB, error) {
	db, err := sql.Open(dbType, "./"+dbName+".db"+enabledFlag)
	return db, err
}

// execute operation into DB based on the query input. For example, performs a CREATE, INSERT, and other
// SQL operation.
func ActionOnTable(database *sql.DB, query string) error {
	statement, err := database.Prepare(query)
	if err != nil {
		fmt.Println(err)
		return err
	}
	_, err = statement.Exec()
	return err
}

//execute a query to retrieve data from DB
func RetrieveData(database *sql.DB, query string) (*sql.Rows, error) {
	results, err := database.Query(query)
	return results, err
}
