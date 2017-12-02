// +build test

package main

import (
	"github.com/ChimeraCoder/anaconda"
	"log"
	"net/url"
	"os"
)

func main() {
	twitterKey, err := ParseTwitterKey(os.Args[3])
	if err != nil {
		log.Fatal(err)
	}

	anaconda.SetConsumerKey(twitterKey.ConsumerKey)
	anaconda.SetConsumerSecret(twitterKey.ConsumerSecret)
	api := anaconda.NewTwitterApi(twitterKey.AccessToken, twitterKey.AccessSecret)

	botConfig, err := ParseBotConfig(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	result, err := RunCmd(os.Args[2], botConfig)
	if err != nil {
		log.Fatal(err)
	}

	result = MakeTweetable(result)
	v := url.Values{}
	v.Set("tweet_mode", "extended")
	_, err = api.PostTweet(result, v)
	log.Println(err)

}
