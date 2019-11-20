package models

import (
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

var DB *sqlx.DB

var schema = `
CREATE TABLE IF NOT EXISTS users (	
    email text,
    slug text
);

CREATE TABLE IF NOT EXISTS presentations (	
    slug text,
	file text,
	author text,
    name text,
    description text,
	json text,
	draft int DEFAULT 1, 
	pages int DEFAULT 0,
    date date
)`

func checkSchema() {
	DB.MustExec(schema)

	if _, err := os.Stat("static/presentation"); os.IsNotExist(err) {
		os.Mkdir("static/presentation", 0755)
	}
}

func Connect() {
	// DB, _ = sqlx.Open("sqlite3", ":memory:")
	DB = sqlx.MustConnect("sqlite3", "db.sqlite") //
	checkSchema()
}
