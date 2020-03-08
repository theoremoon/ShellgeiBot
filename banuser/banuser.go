package main

import (
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"

	"github.com/ChimeraCoder/anaconda"
	_ "github.com/mattn/go-sqlite3"
)

type twitterKeys struct {
	ConsumerKey    string `json:"ConsumerKey"`
	ConsumerSecret string `json:"ConsumerSecret"`
	AccessToken    string `json:"AccessToken"`
	AccessSecret   string `json:"AccessSecret"`
}

func parseTwitterKey(file string) (twitterKeys, error) {
	var k twitterKeys
	raw, err := ioutil.ReadFile(file)
	if err != nil {
		return k, err
	}
	err = json.Unmarshal(raw, &k)
	if err != nil {
		return k, err
	}
	return k, nil
}

// parse flags
func parseFlags() (twitter string, screenName string, err error) {
	t := flag.String("twitter", "", "path to twitter config json file")
	u := flag.String("user", "", "twitter user screen name")
	flag.Parse()
	if t == nil || *t == "" {
		flag.PrintDefaults()
		return "", "", fmt.Errorf("'twitter' option is not specified")
	}
	if u == nil || *u == "" {
		flag.PrintDefaults()
		return "", "", fmt.Errorf("'user' option is not specified")
	}
	return *t, *u, nil
}

func run() error {
	// parse flags
	twitter, screenName, err := parseFlags()
	if err != nil {
		log.Println(err)
		return nil
	}

	// initialize Anaconda
	twitterKey, err := parseTwitterKey(twitter)
	if err != nil {
		return err
	}
	anaconda.SetConsumerKey(twitterKey.ConsumerKey)
	anaconda.SetConsumerSecret(twitterKey.ConsumerSecret)
	api := anaconda.NewTwitterApi(twitterKey.AccessToken, twitterKey.AccessSecret)

	// connect to DB
	db, err := sql.Open("sqlite3", "./database.db")
	if err != nil {
		return err
	}
	defer db.Close()

	// main logic
	self, err := api.GetSelf(url.Values{})
	if err != nil {
		return err
	}
	user, err := api.GetUsersShow(screenName, url.Values{})
	if err != nil {
		return err
	}
	_, err = api.UnblockUserId(user.Id, url.Values{}) // temporarily unblock user to search tweet
	if err != nil {
		return err
	}
	rows, err := db.Query("SELECT tweet_id FROM shellgeis WHERE user_id = ?", user.Id)
	if err != nil {
		return err
	}
	for rows.Next() {
		var tweetID string
		rows.Scan(&tweetID)

		v := url.Values{}
		result, err := api.GetSearch(tweetID, v)
		if err != nil {
			log.Println(err)
			continue
		}

		for _, tweet := range result.Statuses {
			if tweet.User.Id != self.Id {
				continue
			}
			_, err = api.DeleteTweet(tweet.Id, false)
			if err != nil {
				log.Println(err)
				continue
			}
			log.Printf("[+] Deleted id: %d\n", tweet.Id)
		}
	}
	_, err = api.BlockUserId(user.Id, url.Values{}) // re-block user
	if err != nil {
		return err
	}
	return nil
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}

}
