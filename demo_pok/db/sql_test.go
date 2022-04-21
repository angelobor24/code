package db

import (
	"testing"

	"gotest.tools/assert"
)

func TestInitializeDB(t *testing.T) {
	db, err := InitializeDB("sqlite3", "dummy", "")
	assert.Equal(t, db != nil, true)
	assert.Equal(t, err, nil)
}

func TestActionOnTable(t *testing.T) {
	db, _ := InitializeDB("sqlite3", "dummy", "")
	err := ActionOnTable(db, "SELECT * from dummy")
	assert.Equal(t, err != nil, true)
	err = ActionOnTable(db, "CREATE TABLE IF NOT EXISTS dummy (id INTEGER PRIMARY KEY)")
	assert.Equal(t, nil, err)
	err = ActionOnTable(db, "SELECT * from dummy")
	assert.Equal(t, nil, err)
	ActionOnTable(db, "DROP TABLE dummy")
}

func TestInsertTable(t *testing.T) {
	db, _ := InitializeDB("sqlite3", "dummy", "")
	err := ActionOnTable(db, "INSERT INTO dummy VALUES 24")
	assert.Equal(t, err != nil, true)
	ActionOnTable(db, "CREATE TABLE IF NOT EXISTS dummy (id INTEGER PRIMARY KEY)")
	err = ActionOnTable(db, "INSERT INTO dummy VALUES (24)")
	assert.Equal(t, err, nil)
	ActionOnTable(db, "DROP TABLE dummy")
}

func TestRetrieveData(t *testing.T) {
	db, _ := InitializeDB("sqlite3", "dummy", "")
	_, err := RetrieveData(db, "SELECT * FROM dummy")
	assert.Equal(t, err != nil, true)
	ActionOnTable(db, "CREATE TABLE IF NOT EXISTS dummy (id INTEGER PRIMARY KEY)")
	ActionOnTable(db, "INSERT INTO dummy VALUES (24)")
	_, err = RetrieveData(db, "SELECT * FROM dummy")
	assert.Equal(t, err, nil)
	ActionOnTable(db, "DROP TABLE dummy")
}
