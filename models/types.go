package models

import (
	"database/sql"
)

type Entry struct {
	Sort        int    `json:"sort" db:"sort"`
	Audio       string `json:"audio" db:"audio"`
	AudioLength int    `json:"audioLength" db:"audioLength"`
	Image       string `json:"image" db:"image"`
}

type Presentation struct {
	File        string         `json:"file" db:"file"`
	Slug        string         `json:"slug" db:"slug"`
	Author      string         `json:"author" db:"author"`
	Date        sql.NullString `json:"date" db:"date"`
	Name        string         `json:"name" db:"name"`
	Description string         `json:"description" db:"description"`
	Json        sql.NullString `json:"json" db:"json"`
	Pages       int            `json:"pages" db:"pages"`
	Draft       bool           `json:"draft" db:"draft"`
	Entries     []Entry        `json:"entries" db:"entries"`
}
