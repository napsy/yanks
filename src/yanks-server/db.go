package main

import (
	"database/sql"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Db interface {
	getMemStats(app app, from, to time.Time) ([]memEntry, error)
	putMemStats(app app, stats []memEntry) error
}

type sqlDb struct {
	db *sql.DB
}

func (db *sqlDb) getMemStats(app app, from, to time.Time) ([]memEntry, error) {
	rows, err := db.db.Query("SELECT * FROM userinfo")
	if err != nil {
		return nil, err
	}

	entries := []memEntry{}
	for rows.Next() {
		e := memEntry{}
		err = rows.Scan(&e.flat, &e.flatP, &e.sum, &e.cum, &e.cumP, &e.fn)
		if err != nil {
			break
		}
		entries = append(entries, e)
	}
	return entries, nil
}

func (db *sqlDb) putMemStats(app app, stats []memEntry) error {
	_, err := db.db.Prepare("INSERT INTO memstats(username, departname, created) values(?,?,?)")
	return err
}

func (db *sqlDb) prepare() error {
	_, err := db.db.Prepare(`CREATE TABLE 'memstats' (
    'uid' INTEGER PRIMARY KEY AUTOINCREMENT,
    'username' VARCHAR(64) NULL,
    'departname' VARCHAR(64) NULL,
    'created' DATE NULL;`)
	if err != nil {
		return err
	}
	return nil
}

func newSqlDb(filename string) (Db, error) {
	db, err := sql.Open("sqlite3", "./foo.db")
	if err != nil {
		return nil, err
	}
	sql := sqlDb{
		db: db,
	}
	return &sql, nil
}
