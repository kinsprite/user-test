package main

import (
	"log"
	"os"

	"github.com/jmoiron/sqlx"
)

const schema = `
CREATE TABLE user (
    id int AUTO_INCREMENT PRIMARY KEY,
    name text,
    email text
);

CREATE TABLE place (
    country text,
    city text NULL,
    telcode integer
);

ALTER TABLE user AUTO_INCREMENT = 10001;`

var db *sqlx.DB

func initDB() {
	driverName := os.Getenv("SQL_DRIVER_NAME")
	dataSourceName := os.Getenv("SQL_DATA_SOURCE_NAME")

	// MYSQL DSN format: username:password@protocol(address)/dbname?param=value
	var err error
	db, err = sqlx.Open(driverName, dataSourceName)

	if err != nil {
		log.Println("ERROR    Fail to open to the USER DB, ", err)
	}

	err = db.Ping()

	if err != nil {
		log.Println("ERROR    Fail to ping to the USER DB, ", err)
		return
	}

	_, err = db.Exec(schema)

	if err != nil {
		log.Println("ERROR    Fail to create schema to the USER DB, ", err)
		return
	}

	log.Println("INFO    Success to create schema to the USER DB")
}

func closeDB() {
	if db != nil {
		db.Close()
	}
}

func createUserInfoToDB(userInfo *UserInfo) {
	sql := "INSERT INTO user (name, email) VALUES (:name, :email)"

	// if userInfo.ID != 0 {
	// 	sql = "INSERT INTO user (id, name, email) VALUES (:id, :name, :email)"
	// }

	res, err := db.NamedExec(sql, userInfo)

	if err != nil {
		log.Println("ERROR   Insert user info to DB, user name: ", userInfo.Name)
		return
	}

	// if userInfo.ID != 0 {
	// 	return
	// }

	newID, err := res.LastInsertId()

	if err != nil {
		log.Println("INFO   Insert user info to DB with user ID, and can't get the auto inserted ID, user name: ", userInfo.Name)
		return
	}

	userInfo.ID = int(newID)
}

func getUserInfoFromDB(userID int) *UserInfo {
	userInfo := UserInfo{}
	err := db.Get(&userInfo, "SELECT (id, name, email) FROM user WHERE id=$1", userID)

	if err != nil {
		log.Println("FAIL    Get user info from DB, user ID: ", userID)
		return nil
	}

	return &userInfo
}
