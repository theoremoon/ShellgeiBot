package main

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

const schema = `
create table if not exists shellgeis (
  user_id integer,
  screen_name text,
  tweet_id integer,
  shellgei text,
  result text default "",
  error text default "",
  timestamp integer
);
`

func insertResult(db *sql.DB, tweetID int64, result string, err error) error {
	errStr := ""
	if err != nil {
		errStr = err.Error()
	}
	_, err2 := db.Exec("update shellgeis set result=?, error=? where tweet_id=?", result, errStr, tweetID)
	return err2
}

func insertShellGei(db *sql.DB, userID int64, screenName string, tweetID int64, shellgei string, timestamp int64) error {
	_, err := db.Exec("insert into shellgeis(user_id, screen_name, tweet_id, shellgei, timestamp) values (?,?,?,?,?)", userID, screenName, tweetID, shellgei, timestamp)
	if err != nil {
		return err
	}
	return nil
}
