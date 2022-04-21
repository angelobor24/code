package db

import (
	"testing"

	"gotest.tools/assert"
)

func TestNewStorageImpl(t *testing.T) {
	storageImpl := NewStorageImpl("dummy", "", "?h")
	assert.Equal(t, storageImpl.GetStatus(), true)
}
func TestInitializeTables(t *testing.T) {
	var listTables []string
	listTables = append(listTables, "dummy")
	newStorage := NewStorageImpl("sqlite3", "test", "?_foreign_keys=true")
	err := newStorage.InitializeTables(listTables)
	assert.Equal(t, err != nil, true)
	listTables[0] = CreateTrainer
	err = newStorage.InitializeTables(listTables)
	assert.Equal(t, nil, err)
	err = ActionOnTable(DbInstance, "DROP TABLE trainer")
	assert.Equal(t, nil, err)
}

func TestCheckTrainerMt(t *testing.T) {
	newStorage := NewStorageImpl("sqlite3", "test", "?_foreign_keys=true")
	ActionOnTable(DbInstance, CreateTrainer)
	_, err := newStorage.CheckTrainerMt(2)
	assert.Equal(t, err != nil, true)
	ActionOnTable(DbInstance, "INSERT INTO trainer VALUES (24,'dummy','dummy',2)")
	_, err = newStorage.CheckTrainerMt(2)
	assert.Equal(t, err, nil)
	err = ActionOnTable(DbInstance, "DROP TABLE trainer")
	assert.Equal(t, nil, err)

}
