// +build !test

package main

import (
	"bytes"
	"database/sql"
	"encoding/base64"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"strings"

	"github.com/ChimeraCoder/anaconda"
	"github.com/mattn/go-sixel"
	_ "github.com/mattn/go-sqlite3"
)

func processTweet(tweet anaconda.Tweet, self anaconda.User, api *anaconda.TwitterApi, db *sql.DB, config botConfig) {
	// check if it is valid shellgei tweet
	if tweet.RetweetedStatus != nil {
		return
	}
	if !isShellGeiTweet(tweet, config.Tags) {
		return
	}
	if self.Id == tweet.User.Id {
		return
	}
	if !isFollower(api, tweet) {
		return
	}

	t, err := tweet.CreatedAtTime()
	if err != nil {
		log.Println(err)
		return
	}
	text, mediaUrls, err := extractShellgei(tweet, self, api, config.Tags,0)
	if err != nil {
		log.Println(err)
		return
	}

	insertShellGei(db, tweet.User.Id, tweet.User.ScreenName, tweet.Id, text, t.Unix())

	result, b64imgs, err := runCmd(text, mediaUrls, config)
	result = makeTweetable(result)
	insertResult(db, tweet.Id, result, err)

	if err != nil {
		if err.(*stdError) == nil {
			_, _ = api.PostTweet("@theoldmoon0602 internal error", url.Values{})
		}
		return
	}

	if len(result) == 0 && len(b64imgs) == 0 {
		return
	}

	err = tweetResult(api, tweet, result, b64imgs)
	if err != nil {
		log.Println(err)
	}
	return
}

/// ShellgeiBot main function
func botMain(twitterConfigFile, botConfigFile string) {
	twitterKey, err := parseTwitterKey(twitterConfigFile)
	if err != nil {
		log.Fatal(err)
	}

	db, err := sql.Open("sqlite3", "./database.db")
	if err != nil {
		log.Fatal(err)
	}
	_, _ = db.Exec(schema)

	anaconda.SetConsumerKey(twitterKey.ConsumerKey)
	anaconda.SetConsumerSecret(twitterKey.ConsumerSecret)
	api := anaconda.NewTwitterApi(twitterKey.AccessToken, twitterKey.AccessSecret)

	v := url.Values{}
	self, err := api.GetSelf(v)
	if err != nil {
		log.Fatal(err)
	}

	config, err := parseBotConfig(botConfigFile)
	if err != nil {
		log.Fatal(err)
	}
	v.Set("track", strings.Join(config.Tags, ","))
	stream := api.PublicStreamFilter(v)

	for {
		t := <-stream.C
		switch tweet := t.(type) {
		case anaconda.Tweet:
			config, err = parseBotConfig(botConfigFile)
			if err != nil {
				_, _ = api.PostTweet("@theoldmoon0602 Internal error", v)
				log.Fatal(err)
			}

			go func() {
				processTweet(tweet, self, api, db, config)
			}()
		}
	}
}

func botTest(botConfigFile, scriptFile string) {
	config, err := parseBotConfig(botConfigFile)
	if err != nil {
		log.Fatal(err)
	}

	script, err := ioutil.ReadFile(scriptFile)
	if err != nil {
		log.Fatal(err)
	}

	result, b64imgs, err := runCmd(string(script), []string{}, config)
	result = makeTweetable(result)

	if err != nil {
		if err.(*stdError) == nil {
			log.Fatal("internal Error")
		}
		return
	}

	if len(result) == 0 && len(b64imgs) == 0 {
		fmt.Println("No result")
		return
	}

	fmt.Println(result)
	fmt.Println(len(b64imgs))
	for _, b64img := range b64imgs {
		imgBytes, err := base64.StdEncoding.DecodeString(b64img)
		if err != nil {
			log.Println(err)
			continue
		}

		img, _, err := image.Decode(bytes.NewReader(imgBytes))
		if err != nil {
			log.Println(err)
			continue
		}

		sixel.NewEncoder(os.Stdout).Encode(img)
	}
}

func main() {
	if len(os.Args) < 3 {
		log.Fatalf("<Usage>%s: TwitterConfig.json ShellgeiConfig.json | -test ShellgeiConfig.json script", os.Args[0])
	}

	if os.Args[1] == "-test" {
		// testing mode
		botTest(os.Args[2], os.Args[3])
	} else {
		// normal mode
		botMain(os.Args[1], os.Args[2])
	}
}
