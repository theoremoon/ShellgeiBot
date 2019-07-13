// +build !test

package main

import (
	"bytes"
	"database/sql"
	"encoding/base64"
	"fmt"
	"github.com/ChimeraCoder/anaconda"
	"github.com/mattn/go-sixel"
	_ "github.com/mattn/go-sqlite3"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"strings"
)

func ProcessTweet(tweet anaconda.Tweet, self anaconda.User, api *anaconda.TwitterApi, db *sql.DB, botConfig BotConfig) {
	// check if it is valid shellgei tweet
	if tweet.RetweetedStatus != nil {
		return
	}
	if !IsShellGeiTweet(tweet.FullText, botConfig.Tags) {
		return
	}
	if self.Id == tweet.User.Id {
		return
	}
	if !IsFollower(api, tweet) {
		return
	}

	t, err := tweet.CreatedAtTime()
	if err != nil {
		log.Println(err)
		return
	}
	text, media_urls, err := ExtractShellgei(tweet, self, api, botConfig.Tags)
	if err != nil {
		log.Println(err)
		return
	}

	InsertShellGei(db, tweet.User.Id, tweet.User.ScreenName, tweet.Id, text, t.Unix())

	result, b64imgs, err := RunCmd(text, media_urls, botConfig)
	result = MakeTweetable(result)
	InsertResult(db, tweet.Id, result, err)

	if err != nil {
		if err.(*StdError) == nil {
			_, _ = api.PostTweet("@theoldmoon0602 internal error", url.Values{})
		}
		return
	}

	if len(result) == 0 && len(b64imgs) == 0 {
		return
	}

	err = TweetResult(api, tweet, result, b64imgs)
	if err != nil {
		log.Println(err)
	}
	return
}

/// ShellgeiBot main function
func botMain(twitterConfigFile, botConfigFile string) {
	twitterKey, err := ParseTwitterKey(twitterConfigFile)
	if err != nil {
		log.Fatal(err)
	}

	db, err := sql.Open("sqlite3", "./database.db")
	if err != nil {
		log.Fatal(err)
	}
	_, _ = db.Exec(Schema)

	anaconda.SetConsumerKey(twitterKey.ConsumerKey)
	anaconda.SetConsumerSecret(twitterKey.ConsumerSecret)
	api := anaconda.NewTwitterApi(twitterKey.AccessToken, twitterKey.AccessSecret)

	v := url.Values{}
	self, err := api.GetSelf(v)
	if err != nil {
		log.Fatal(err)
	}

	botConfig, err := ParseBotConfig(botConfigFile)
	if err != nil {
		log.Fatal(err)
	}
	v.Set("track", strings.Join(botConfig.Tags, ","))
	stream := api.PublicStreamFilter(v)

	for {
		t := <-stream.C
		switch tweet := t.(type) {
		case anaconda.Tweet:
			botConfig, err = ParseBotConfig(botConfigFile)
			if err != nil {
				_, _ = api.PostTweet("@theoldmoon0602 Internal error", v)
				log.Fatal(err)
			}

			go func() {
				ProcessTweet(tweet, self, api, db, botConfig)
			}()
		}
	}
}

func botTest(botConfigFile, scriptFile string) {
	botConfig, err := ParseBotConfig(botConfigFile)
	if err != nil {
		log.Fatal(err)
	}

	script, err := ioutil.ReadFile(scriptFile)
	if err != nil {
		log.Fatal(err)
	}

	result, b64imgs, err := RunCmd(string(script), []string{}, botConfig)
	result = MakeTweetable(result)

	if err != nil {
		if err.(*StdError) == nil {
			log.Fatal("internal Error")
		}
		return
	}

	if len(result) == 0 && len(b64imgs) == 0 {
		fmt.Println("No result")
		return
	}

	fmt.Println(result)
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
