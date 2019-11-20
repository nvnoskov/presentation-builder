package models

import (
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
	author text NOT NULL,
    name text NOT NULL,
    description text NOT NULL,
	json text,
	draft int DEFAULT 1, 
	pages int DEFAULT 0,
    date date
)`

func checkSchema() {
	DB.MustExec(schema)
}

func Connect() {
	// DB, _ = sqlx.Open("sqlite3", ":memory:")
	DB = sqlx.MustConnect("sqlite3", "db.sqlite") //
	checkSchema()
}
