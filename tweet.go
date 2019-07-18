package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html"
	"io/ioutil"
	"log"
	"net/url"
	"strings"
	"unicode"

	"github.com/ChimeraCoder/anaconda"
)

type TwitterKeys struct {
	ConsumerKey    string `json:"ConsumerKey"`
	ConsumerSecret string `json:"ConsumerSecret"`
	AccessToken    string `json:"AccessToken"`
	AccessSecret   string `json:"AccessSecret"`
}

// Twitter.Entities.Hashtagsが無名構造体のため、該当構造体の初期化に手間がかかる
// 。全く同じ構造体ではあるが、実装を容易にするため部分的に切り出して型を定義。
type TweetEntitiesHashtags []struct {
	Indices []int
	Text    string
}

func ExtractShellgei(tweet anaconda.Tweet, self anaconda.User, api *anaconda.TwitterApi, tags []string) (string, []string, error) {
	// self recursion
	if tweet.QuotedStatusID == tweet.Id {
		return "", nil, fmt.Errorf("self recursion")
	}

	// if it is quoted tweet of shellgeibot's tweet
	// then dig deeper (ignore shellgeibot's output)
	if tweet.User.Id == self.Id {
		if tweet.QuotedStatusID == 0 { // will never be true
			return "", nil, fmt.Errorf("quote tweet by ownself")
		}

		// get quoted tweet and dig deeper
		v := url.Values{}
		quoted, err := api.GetTweet(tweet.QuotedStatusID, v)
		if err != nil {
			return "", nil, err
		}
		return ExtractShellgei(quoted, self, api, tags)
	}

	// get tweet text
	text := tweet.FullText

	// expand url
	for _, url := range tweet.Entities.Urls {
		if strings.HasPrefix(url.Expanded_url, "https://") {
			text = strings.Replace(text, url.Url, url.Expanded_url[len("https://"):], -1)
		} else if strings.HasPrefix(url.Expanded_url, "http://") {
			text = strings.Replace(text, url.Url, url.Expanded_url[len("http://"):], -1)
		}
	}
	// list of picture url
	media_urls := make([]string, 0, 4)
	for _, media := range tweet.ExtendedEntities.Media {
		media_urls = append(media_urls, media.Media_url_https)

		// remove media url
		text = strings.Replace(text, media.Url, "", -1)
	}

	// treat
	text = html.UnescapeString(text)
	text = RemoveMentionSymbol(self, text)

	// remove tags
	shellgei := removeTags(text, tweet.Entities.Hashtags, tags)

	if tweet.QuotedStatusID == 0 {
		return shellgei, media_urls, nil
	}

	// tweet chain
	v := url.Values{}
	quoted, err := api.GetTweet(tweet.QuotedStatusID, v)
	if err != nil {
		return "", nil, err
	}

	quote_text, quote_urls, err := ExtractShellgei(quoted, self, api, tags)
	if err != nil {
		return "", nil, err
	}
	return quote_text + shellgei, append(quote_urls, media_urls...), nil
}

func removeTags(text string, hashtags TweetEntitiesHashtags, tags []string) string {
	const removeMark = rune(0xFFFE)

	rtext := []rune(text)
	for _, tag := range hashtags {
		for _, t := range tags {
			if tag.Text == t {
				if len(tag.Indices) < 2 {
					msg := fmt.Sprintf("[WARN] tag indices was until 2. text = %s, rtext = %s, tag indices length = %d, hashtags = %v, tag = %s", text, string(rtext), len(tag.Indices), hashtags, t)
					log.Println(msg)
					continue
				}

				tagStartPos := tag.Indices[0]
				tagEndPos := tag.Indices[1]

				// Set remove marks to tag range.
				for i := tagStartPos; i < tagEndPos; i++ {
					rtext[i] = removeMark
				}
			}
		}
	}
	var b bytes.Buffer
	for _, v := range rtext {
		if v != removeMark {
			b.WriteRune(v)
		}
	}
	return strings.TrimSpace(b.String())
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
		if a[i] <= unicode.MaxASCII || unicode.Is(unicode.Latin, a[i]) {
			l++
		} else {
			l += 2
		}
		if l > 280 {
			break
		}
	}
	return string(a[:i])
}

func TweetUrl(tweet anaconda.Tweet) string {
	return "https://twitter.com/" + tweet.User.ScreenName + "/status/" + tweet.IdStr
}

func TweetResult(api *anaconda.TwitterApi, tweet anaconda.Tweet, result string, b64imgs []string) error {
	v := url.Values{}
	media_ids := make([]string, 0, 4)

	// with image
	for _, b64img := range b64imgs {
		media, err := api.UploadMedia(b64img)
		if err != nil {
			log.Println(err)
		} else {
			media_ids = append(media_ids, media.MediaIDString)
		}
	}
	v.Add("media_ids", strings.Join(media_ids, ","))

	if len(b64imgs) == 0 {
		v.Set("tweet_mode", "extended")
		v.Set("attachment_url", TweetUrl(tweet))
	} else {
		result = result + " " + TweetUrl(tweet)
	}

	// post message
	_, err := api.PostTweet(result, v)
	return err
}

func IsShellGeiTweet(tweet anaconda.Tweet, tags []string) bool {
	for _, tag := range tweet.Entities.Hashtags {
		for _, t := range tags {
			if t == tag.Text {
				return true
			}
		}
	}
	return false
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
