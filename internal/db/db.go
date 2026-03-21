package db

import (
	"database/sql"
	"log"
)

func DBInit() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "./envs.db")
	if err != nil {
		return nil, err
	}

	sqlStmt := `
    CREATE TABLE IF NOT EXISTS project (
        id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
        name TEXT,
	    	created_at TEXT DEFAULT CURRENT_TIMESTAMP
    );
    `
	_, err = db.Exec(sqlStmt)
	if err != nil {
		return nil, err
	}

	sqlStmt = `
    CREATE TABLE IF NOT EXISTS project (
        id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
        name TEXT,
	    	created_at TEXT DEFAULT CURRENT_TIMESTAMP
    );
    `

	_, err = db.Exec(sqlStmt)
	if err != nil {
		return nil, err
	}
	log.Println("Table 'project' created successfully")
	return db, nil
}
