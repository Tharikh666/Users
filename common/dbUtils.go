package common

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

var DBConnection *sql.DB

var DBConnectionEstablished bool = false

func GetDBConnection() (*sql.DB, bool) {
	if DBConnectionEstablished {
		return DBConnection, true
	} else {
		err := OpenDBConnection()
		if err != nil {
			log.Fatalln(err)
			return DBConnection, false
		} else {
			return DBConnection, true
		}
	}
}

func OpenDBConnection() error {
	connectionString := Config.ConnectionString

	myDB, err := sql.Open("postgres", connectionString)

	if err != nil {
		DBConnectionEstablished = false
		DBConnection = nil
		log.Fatalln("opening DB Connection failed")
		return err
	} else {
		DBConnectionEstablished = true
		DBConnection = myDB
		return nil
	}
}
