package db

import (
	"database/sql"
	"fmt"
	"strconv"
)

var (
	ListTables = [...]string{CreateTrainer, CreateQuote, CreateBaseQuote, CreatePayedQuote, CreateExtraQuote}
)

const (
	CreateTrainer               = "CREATE TABLE IF NOT EXISTS trainer (id_trainer INTEGER PRIMARY KEY, name TEXT NOT NULL, surname TEXT NOT NULL, mt INTEGER NOT NULL)"
	CreateQuote                 = "CREATE TABLE IF NOT EXISTS quote (id_trainer INTEGER, pokemon TEXT, price REAL NOT NULL, PRIMARY KEY (id_trainer, pokemon),CONSTRAINT f1 FOREIGN KEY (id_trainer) REFERENCES trainer (id_trainer) ON DELETE CASCADE)"
	CreateBaseQuote             = "CREATE TABLE IF NOT EXISTS basequote (category TEXT PRIMARY KEY, price REAL NOT NULL)"
	CreateExtraQuote            = "CREATE TABLE IF NOT EXISTS extraquote (category TEXT PRIMARY KEY, price REAL NOT NULL)"
	CreatePayedQuote            = "CREATE TABLE IF NOT EXISTS payedquote (id INTEGER PRIMARY KEY, id_trainer INTEGER, pokemon TEXT NOT NULL, price INTEGER NOT NULL, timestamp DATETIME NOT NULL, CONSTRAINT f2 FOREIGN KEY (id_trainer) REFERENCES trainer (id_trainer) ON DELETE CASCADE)"
	insertTrainer               = "INSERT INTO trainer (id_trainer,name,surname, mt) VALUES"
	insertQuote                 = "INSERT INTO quote (id_trainer,pokemon,price) VALUES"
	insertBaseQuote             = "INSERT INTO basequote (category, price) VALUES"
	insertExtraQuote            = "INSERT INTO extraquote (category, price) VALUES"
	insertPayedQuote            = "INSERT INTO payedquote (id,id_trainer,pokemon,price,timestamp) VALUES"
	getFromTrainer              = "SELECT id_trainer, name, surname, mt FROM trainer"
	getFromQuote                = "SELECT id_trainer, pokemon, price FROM quote"
	getFromQuoteFiltered        = "SELECT price FROM quote WHERE id_trainer=%s AND pokemon='%s'"
	getFromBaseQuote            = "SELECT price FROM basequote WHERE category="
	getFromExtraQuote           = "SELECT price FROM extraquote WHERE category="
	checkIfTrainerExistFilterMt = "SELECT id_trainer FROM trainer WHERE mt="
	checkIfTrainerExistFilterId = "SELECT id_trainer FROM trainer WHERE id_trainer="
	deleteFromQuote             = "DELETE FROM quote WHERE id_trainer=%s AND pokemon='%s'"
)

var (
	DbInstance *sql.DB
)

type Storage interface {
	GetStatus() bool
	InitializeTables(listTables []string) error
	InsertTrainer(id_trainer int, name string, surname string, mt int) error
	InsertBaseQuote(category string, price float32) error
	InsertExtraQuote(category string, price float32) error
	InsertQuote(id_trainer int, pokemon string, price float32) error
	InsertPayedQuote(id int64, id_trainer int, pokemon string, price float32) error
	DeleteFromQuote(id_trainer int, pokemon string) error
	GetFromBaseQuote(category string) (*sql.Rows, error)
	GetFromExtraQuote(category string) (*sql.Rows, error)
	GetFromQuoteFiltered(id_trainer int, pokemon string) (*sql.Rows, error)
	CheckTrainerMt(mt int) (int, error)
	CheckTrainerExist(id int) (*sql.Rows, error)
}

type StorageImpl struct {
	status bool
}

// instance of struct that implements Storage interface
// This struct is responsible to handle the DB, so initialize also the DB structure,
func NewStorageImpl(dbType string, nameDB string, option string) Storage {
	storageImpl := StorageImpl{}
	err := createDB(dbType, nameDB, option)
	if err != nil {
		fmt.Println(err)
		storageImpl.status = false
	}
	storageImpl.status = true
	return &storageImpl
}

// return the DB status. If false, the DB instantiation is failed
func (storageImpl *StorageImpl) GetStatus() bool {
	return storageImpl.status
}

// layer to call the DB structure initialization
func createDB(dbType string, nameDB string, option string) error {
	var err error
	DbInstance, err = InitializeDB(dbType, nameDB, option)
	return err
}

// taking a list of query, performs n CREATE operation, with n=len(listBables)
func (storageImpl *StorageImpl) InitializeTables(listTables []string) error {
	for _, value := range listTables {
		err := ActionOnTable(DbInstance, value)
		if err != nil {
			fmt.Println(err)
			return err
		}
	}
	return nil
}

// operation to Insert into trainer table
func (storageImpl *StorageImpl) InsertTrainer(id_trainer int, name string, surname string, mt int) error {
	return ActionOnTable(DbInstance, insertTrainer+"("+strconv.Itoa(id_trainer)+",'"+name+"','"+surname+"',"+strconv.Itoa(mt)+")")
}

// operation to Insert into basequote table
func (storageImpl *StorageImpl) InsertBaseQuote(category string, price float32) error {
	return ActionOnTable(DbInstance, insertBaseQuote+"('"+category+"',"+fmt.Sprint(price)+")")
}

// operation to Insert into extraquote table
func (storageImpl *StorageImpl) InsertExtraQuote(category string, price float32) error {
	return ActionOnTable(DbInstance, insertExtraQuote+"('"+category+"',"+fmt.Sprint(price)+")")
}

// operation to Insert into quote table
func (storageImpl *StorageImpl) InsertQuote(id_trainer int, pokemon string, price float32) error {
	return ActionOnTable(DbInstance, insertQuote+"("+strconv.Itoa(id_trainer)+",'"+pokemon+"',"+fmt.Sprint(price)+")")
}

// operation to Insert into payedquote table
func (storageImpl *StorageImpl) InsertPayedQuote(id int64, id_trainer int, pokemon string, price float32) error {
	return ActionOnTable(DbInstance, insertPayedQuote+"("+strconv.FormatInt(id, 10)+","+strconv.Itoa(id_trainer)+",'"+pokemon+"',"+fmt.Sprint(price)+",CURRENT_TIMESTAMP)")
}

// operation to Delete from quote table
func (storageImpl *StorageImpl) DeleteFromQuote(id_trainer int, pokemon string) error {
	return ActionOnTable(DbInstance, fmt.Sprintf(deleteFromQuote, strconv.Itoa(id_trainer), pokemon))
}

// operation to retrieve from basequote table
func (storageImpl *StorageImpl) GetFromBaseQuote(category string) (*sql.Rows, error) {
	return RetrieveData(DbInstance, getFromBaseQuote+"'"+category+"'")
}

// operation to retrieve from extraquote table
func (storageImpl *StorageImpl) GetFromExtraQuote(category string) (*sql.Rows, error) {
	return RetrieveData(DbInstance, getFromExtraQuote+"'"+category+"'")
}

// operation to retrieve from quote table
func (storageImpl *StorageImpl) GetFromQuoteFiltered(id_trainer int, pokemon string) (*sql.Rows, error) {
	return RetrieveData(DbInstance, fmt.Sprintf(getFromQuoteFiltered, strconv.Itoa(id_trainer), pokemon))
}

// check if the MT is valid
func (storageImpl *StorageImpl) CheckTrainerMt(mt int) (int, error) {
	var id int
	result, err := RetrieveData(DbInstance, checkIfTrainerExistFilterMt+strconv.Itoa(mt))
	defer result.Close()
	if result.Next() {
		result.Scan(&id)
	} else {
		return id, sql.ErrNoRows
	}
	return id, err
}

// execute the retrieve of id related to the provided mt, if exist
func (storageImpl *StorageImpl) CheckTrainerExist(id int) (*sql.Rows, error) {
	return RetrieveData(DbInstance, checkIfTrainerExistFilterId+strconv.Itoa(id))
}

// start a transaction
func ReturnTransactionContext() (*sql.Tx, error) {
	return DbInstance.Begin()
}
