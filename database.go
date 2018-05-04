package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

const Schema = `
create table if not exists shellgeis (
	user_id integer,
	screen_name text,
	tweet_id integer,
	shellgei text,
	result text default "",
	error text default "",
	timestamp string
);
`

func InsertResult(db *sql.DB, tweet_id int64, result string, err error) error {
	err_str := ""
	if err != nil {
		err_str = err.Error()
	}
	_, err2 := db.Exec("update shellgeis set result=?, error=? where tweet_id=?", result, err_str, tweet_id)
	return err2
}

func InsertShellGei(db *sql.DB, user_id int64, screen_name string, tweet_id int64, shellgei string, timestamp int64) error {
	_, err := db.Exec("insert into shellgeis(user_id, screen_name, tweet_id, shellgei, timestamp) values (?,?,?,?,?)", user_id, screen_name, tweet_id, shellgei, timestamp)
	if err != nil {
		return err
	}
	return nil
}
