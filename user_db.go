package main

import (
	"log"
	"os"

	"github.com/go-sql-driver/mysql"
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
			"charset":         "utf8",
			"multiStatements": "true",
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

	_, err = db.Exec(schema)

	if err != nil {
		log.Println("ERROR    Fail to create schema to the USER DB, ", err)
		return
	}

	log.Println("INFO    Success to create schema to the USER DB")
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
