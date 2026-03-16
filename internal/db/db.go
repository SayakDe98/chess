package db

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func InitDB() *sql.DB {
	var db *sql.DB

	dsn := `root:@tcp(127.0.0.1:3306)/chess?parseTime=true`

	db, err := sql.Open("mysql", dsn)

	if err != nil {
		log.Fatal("Unable to connect to Database")
	}
	if err := db.Ping(); err != nil {
		log.Fatal("Database unreachable")
	}
	return db
}
