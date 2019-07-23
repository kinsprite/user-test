package main

import (
	"log"
	"os"

	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var schema = `
CREATE TABLE user (
    id int,
    name text,
    email text
);

CREATE TABLE place (
    country text,
    city text NULL,
    telcode integer
)`

// MYSQL DSN format: username:password@protocol(address)/dbname?param=value
var sqlDSN string
var db sqlx.DB

func initDB() {
	config := mysql.Config{
		User:   os.Getenv("MYSQL_USER"),
		Passwd: os.Getenv("MYSQL_PASSWORD"),
		Net:    os.Getenv("MYSQL_NET"),
		Addr:   os.Getenv("MYSQL_ADDR"),
		DBName: os.Getenv("MYSQL_DATABASE"),
		Params: map[string]string{
			"charset": "utf8",
			// "allowOldPasswords": "true"
		},
	}

	sqlDSN = config.FormatDSN()
	log.Println("INFO   SQL's DSN: ", sqlDSN)

	db, err := sqlx.Open("mysql", sqlDSN)

	if err != nil {
		log.Println("ERROR    Fail to open to the USER DB, ", err)
	}

	err = db.Ping()

	if err != nil {
		log.Println("ERROR    Fail to ping to the USER DB, ", err)
		return
	}

	// exec the schema or fail; multi-statement Exec behavior varies between
	// database drivers;  pq will exec them all, sqlite3 won't, ymmv
	db.MustExec(schema)

	log.Println("INFO    Success to create schema to the USER DB")
}

func getUserInfoFromDB() {

}
