package main

import (
	"encoding/json"
	"github.com/ChimeraCoder/anaconda"
	"io/ioutil"
	"net/url"
	"strings"
	"unicode"
)

type TwitterKeys struct {
	ConsumerKey    string `json:"ConsumerKey"`
	ConsumerSecret string `json:"ConsumerSecret"`
	AccessToken    string `json:"AccessToken"`
	AccessSecret   string `json:"AccessSecret"`
}

func ParseTwitterKey(file string) (TwitterKeys, error) {
	var k TwitterKeys
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

func MakeTweetable(text string) string {
	a := []rune(text)
	i := 0
	l := 0
	for ; i < len(a); i++ {
		if unicode.Is(unicode.Latin, a[i]) || unicode.IsSpace(a[i]) {
			l++
		} else {
			l += 2
		}
		if l > 280 {
			i--
			break
		}
	}
	return string(a[:i])
}

func TweetUrl(tweet anaconda.Tweet) string {
	return "https://twitter.com/" + tweet.User.ScreenName + "/status/" + tweet.IdStr
}

func TweetResult(api *anaconda.TwitterApi, tweet anaconda.Tweet, result string) error {
	/// post done message
	v := url.Values{}
	v.Set("tweet_mode", "extended")
	v.Set("attachment_url", TweetUrl(tweet))
	_, err := api.PostTweet(result, v)
	return err
}

func IsShellGeiTweet(tweet string, tags []string) (bool, string) {
	flag := false
	for _, t := range tags {
		if strings.Contains(tweet, t) {
			flag = true
			tweet = strings.Replace(tweet, t, "", -1)
		}
	}
	return flag, tweet
}

func RemoveMentionSymbol(self anaconda.User, tweet string) string {
	return strings.Replace(tweet, "@"+self.ScreenName, "", -1)
}

func IsFollower(api *anaconda.TwitterApi, tweet anaconda.Tweet) bool {
	v := url.Values{}
	u, err := api.GetUsersShowById(tweet.User.Id, v)
	if err != nil {
		return false
	}
	return u.Following
}
