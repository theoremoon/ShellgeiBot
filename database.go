package main

import (
	"database/sql"
	"github.com/ChimeraCoder/anaconda"
	_ "github.com/mattn/go-sqlite3"
)

const Schema = `
create table if not exists shellgeis (
	user_id text,
	tweet_id integer,
	shellgei text,
	timestamp integer
);
`
const Schema2 = `
create table if not exists errors (
	error_text text,
	shellgei text
);
`

func InsertError(db *sql.DB, err error, text string) error {
	_, err2 := db.Exec("insert into errors(error_text, shellgei) values (?,?)", err.Error(), text)
	return err2
}

func InsertShellGei(db *sql.DB, tweet anaconda.Tweet, text string) error {
	now, err := tweet.CreatedAtTime()
	if err != nil {
		return err
	}
	_, err = db.Exec("insert into shellgeis(user_id, tweet_id, shellgei, timestamp) values (?,?,?,?)", tweet.User.IdStr, tweet.Id, text, now.Unix())
	if err != nil {
		return err
	}
	return nil
}
