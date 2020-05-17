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

type twitterKeys struct {
	ConsumerKey    string `json:"ConsumerKey"`
	ConsumerSecret string `json:"ConsumerSecret"`
	AccessToken    string `json:"AccessToken"`
	AccessSecret   string `json:"AccessSecret"`
}

// tweetEntitiesHashtags Twitter.Entities.Hashtagsが無名構造体のため、該当構造
// 体の初期化に手間がかかる。全く同じ構造体ではあるが、実装を容易にするため部
// 分的に切り出して型を定義。
type tweetEntitiesHashtags []struct {
	Indices []int
	Text    string
}

func extractShellgei(tweet anaconda.Tweet, self anaconda.User, api *anaconda.TwitterApi, tags []string, checked []int64) (string, []string, error) {
	// self recursion
	if tweet.QuotedStatusID == tweet.Id {
		return "", nil, fmt.Errorf("self recursion")
	}

	//Recursion Detected
	for _, v := range checked {
		if tweet.QuotedStatusID == v {
			return "", nil, fmt.Errorf("recursion detected")
		}
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
		return extractShellgei(quoted, self, api, tags, append(checked, tweet.Id))
	}

	// get tweet text
	text := tweet.FullText

	// remove trigger tags (do this first to avoid range error)
	text = removeTags(text, tweetEntitiesHashtags(tweet.Entities.Hashtags), tweetEntitiesHashtags(tweet.ExtendedEntities.Hashtags), tags)

	// expand url
	for _, url := range tweet.Entities.Urls {
		if strings.HasPrefix(url.Expanded_url, "https://") {
			text = strings.Replace(text, url.Url, url.Expanded_url[len("https://"):], -1)
		} else if strings.HasPrefix(url.Expanded_url, "http://") {
			text = strings.Replace(text, url.Url, url.Expanded_url[len("http://"):], -1)
		}
	}

	// list of picture url
	mediaUrls := make([]string, 0, 4)
	for _, media := range tweet.ExtendedEntities.Media {
		mediaUrls = append(mediaUrls, media.Media_url_https)

		// remove media url
		text = strings.Replace(text, media.Url, "", -1)
	}

	text = html.UnescapeString(text)
	text = removeMentionSymbol(self, text)

	// remove tags
	shellgei := text

	if tweet.QuotedStatusID == 0 {
		return shellgei, mediaUrls, nil
	}

	// tweet chain
	v := url.Values{}
	quoted, err := api.GetTweet(tweet.QuotedStatusID, v)
	if err != nil {
		return "", nil, err
	}

	quoteText, quoteUrls, err := extractShellgei(quoted, self, api, tags, append(checked, tweet.Id))
	if err != nil {
		return "", nil, err
	}
	return quoteText + shellgei, append(quoteUrls, mediaUrls...), nil
}

func removeTags(text string, hashtags, extHashtags tweetEntitiesHashtags, searchTags []string) string {
	const removeMark = rune(0xFFFE)

	rtext := []rune(text)
	for _, tags := range []tweetEntitiesHashtags{hashtags, extHashtags} {
		for _, tag := range tags {
			for _, searchTag := range searchTags {
				if tag.Text != searchTag {
					continue
				}

				if len(tag.Indices) < 2 {
					log.Printf("[WARN] tag indices < 2. text = %s, tag indices %v, tag = %v\n", text, tag.Indices, tag)
					continue
				}

				// Set remove marks to tag range.
				for i := tag.Indices[0]; i < tag.Indices[1]; i++ {
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

func makeTweetable(text string, untrues []string) string {
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

	result := string(a[:i])

	// big brother is watching your tweet for everyone's happy
	for _, untrue := range untrues {
		if strings.Index(result, untrue) >= 0 {
			return ""
		}
	}
	return result
}

func tweetURL(tweet anaconda.Tweet) string {
	return "https://twitter.com/" + tweet.User.ScreenName + "/status/" + tweet.IdStr
}

func tweetResult(api *anaconda.TwitterApi, tweet anaconda.Tweet, result string, b64imgs []string) error {
	v := url.Values{}
	mediaIds := make([]string, 0, 4)

	// with image
	for _, b64img := range b64imgs {
		media, err := api.UploadMedia(b64img)
		if err != nil {
			log.Println(err)
		} else {
			mediaIds = append(mediaIds, media.MediaIDString)
		}
	}
	v.Add("media_ids", strings.Join(mediaIds, ","))

	if len(b64imgs) == 0 {
		v.Set("tweet_mode", "extended")
		v.Set("attachment_url", tweetURL(tweet))
	} else {
		result = result + " " + tweetURL(tweet)
	}

	// post message
	_, err := api.PostTweet(result, v)
	return err
}

func isShellGeiTweet(tweet anaconda.Tweet, tags []string) bool {
	for _, tag := range tweet.Entities.Hashtags {
		for _, t := range tags {
			if t == tag.Text {
				return true
			}
		}
	}
	return false
}

func removeMentionSymbol(self anaconda.User, tweet string) string {
	return strings.Replace(tweet, "@"+self.ScreenName, "", -1)
}

func isFollower(api *anaconda.TwitterApi, tweet anaconda.Tweet) bool {
	v := url.Values{}
	u, err := api.GetUsersShowById(tweet.User.Id, v)
	if err != nil {
		return false
	}
	return u.Following
}
