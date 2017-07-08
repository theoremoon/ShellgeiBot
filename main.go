package main

import (
	"bytes"
	"crypto/rand"
	"database/sql"
	"github.com/ChimeraCoder/anaconda"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/net/context"
	"html"
	"log"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	Image       string
	WorkDir     string
	TimeoutTime time.Duration
}

func DoShellGei(c Config, script string) (string, error) {
	scriptFileName, err := RandStr(10)
	if err != nil {
		return "", err
	}

	scriptFilePath := filepath.Join(c.WorkDir, scriptFileName)
	scriptFile, err := os.Create(scriptFilePath)
	if err != nil {
		return "", err
	}

	defer scriptFile.Close()
	defer os.Remove(scriptFilePath)

	_, err = scriptFile.WriteString(script)
	if err != nil {
		return "", err
	}

	cmd := exec.Command("docker", "run", "--net=none", "--name", scriptFileName, "-v", scriptFilePath+":/"+scriptFileName, c.Image, "bash", "/"+scriptFileName)
	defer func() {
		cmd := exec.Command("docker", "stop", scriptFileName)
		cmd.Run()
		cmd = exec.Command("docker", "rm", scriptFileName)
		cmd.Run()
	}()

	var out bytes.Buffer
	cmd.Stdout = &out

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, c.TimeoutTime)
	defer cancel()

	errChan := make(chan error, 1)
	go func(ctx context.Context) { errChan <- cmd.Run() }(ctx)

	select {
	case <-ctx.Done():
		err = ctx.Err()
		if err != nil {
			return "", err
		}
	case err = <-errChan:
		if err != nil {
			return "", err
		}
		return out.String(), nil
	}
	return out.String(), nil

}

func RandStr(length int) (string, error) {
	const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	randstr := make([]byte, 0, length)
	keys := make([]byte, length)
	_, err := rand.Read(keys)
	if err != nil {
		return "", err
	}
	for _, v := range keys {
		k := int(v) % len(chars)
		randstr = append(randstr, chars[k])
	}

	return string(randstr), nil
}

func TweetUrl(tweet anaconda.Tweet) string {
	return "https://twitter.com/" + tweet.User.ScreenName + "/status/" + tweet.IdStr
}

func Tweetable(text string) string {
	a := []rune(text)
	if len(a) < 140 {
		return text
	}
	return string(a[:140])
}

func TweetResult(api *anaconda.TwitterApi, tweet anaconda.Tweet, result string) error {
	/// post done message
	v := url.Values{}
	v.Set("attachment_url", TweetUrl(tweet))
	_, err := api.PostTweet(Tweetable(result), v)
	if err != nil {
		return err
	}

	return nil
}

func IsShellGeiTweet(tweet string) (bool, string) {
	tags := []string{"#シェル芸", "#危険シェル芸"}
	for _, t := range tags {
		if strings.Contains(tweet, t) {
			return true, strings.Replace(tweet, t, "", -1)
		}
	}
	return false, ""
}

func RemoveMentionSymbol(self anaconda.User, tweet string) string {
	return strings.Replace(tweet, "@"+self.ScreenName, "", -1)
}

const schema = `
create table if not exists shellgeis (
	user_id text,
	tweet_id integer,
	shellgei text,
	timestamp integer
);
`

func InsertShellGei(db *sql.DB, tweet anaconda.Tweet) error {
	now, err := tweet.CreatedAtTime()
	if err != nil {
		return err
	}
	_, err = db.Exec("insert into shellgeis(user_id, tweet_id, shellgei, timestamp) values (?,?,?,?)", tweet.User.IdStr, tweet.Id, tweet.Text, now.Unix())
	if err != nil {
		return err
	}
	return nil
}

func RecentShellgeiCount(db *sql.DB, timeRange int, tweet anaconda.Tweet) (int, error) {
	now, err := tweet.CreatedAtTime()
	if err != nil {
		return 0, err
	}
	var cnt int
	err = db.QueryRow("select count(*) from shellgeis where user_id=? and timestamp > ?", tweet.User.IdStr, now.Unix()-int64(timeRange)).Scan(&cnt)
	return cnt, err
}

func IsFollower(api *anaconda.TwitterApi, tweet anaconda.Tweet) bool {
	v := url.Values{}
	u, err := api.GetUsersShowById(tweet.User.Id, v)
	if err != nil {
		return false
	}
	return u.Following
}

func main() {
	if len(os.Args) < 8 {
		log.Println("7 arguments required. Consumer key, Consumer secret, Access token, Access token secret, timeout[sec], docker image name shellgei_per_minutes")
		return
	}

	db, err := sql.Open("sqlite3", "./database.db")
	if err != nil {
		log.Println(err)
		return
	}
	db.Exec(schema)

	anaconda.SetConsumerKey(os.Args[1])
	anaconda.SetConsumerSecret(os.Args[2])

	api := anaconda.NewTwitterApi(os.Args[3], os.Args[4])
	timeOutSec, err := strconv.Atoi(os.Args[5])
	if err != nil {
		log.Println(err)
		return
	}
	image := os.Args[6]

	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	config := Config{
		image,
		dir,
		time.Duration(timeOutSec) * time.Second,
	}
	shellgeiPerMinutes, err := strconv.Atoi(os.Args[7])
	if err != nil {
		log.Println(err)
		return
	}

	v := url.Values{}
	self, err := api.GetSelf(v)
	if err != nil {
		log.Println(err)
		return
	}

	stream := api.UserStream(v)

	for {
		t := <-stream.C
		switch tweet := t.(type) {
		case anaconda.Tweet:
			go func() {
				if tweet.RetweetedStatus != nil {
					return
				}
				is, text := IsShellGeiTweet(tweet.Text)
				if !is {
					return
				}
				if self.Id == tweet.User.Id {
					log.Println("self shellgei")
					return
				}
				if !IsFollower(api, tweet) {
					log.Println("not following" + tweet.User.ScreenName)
					return
				}
				text = html.UnescapeString(text)
				text = RemoveMentionSymbol(self, text)
				cnt, err := RecentShellgeiCount(db, 60, tweet)
				InsertShellGei(db, tweet)
				if err != nil {
					log.Println(err)
					return
				}
				if cnt > shellgeiPerMinutes {
					log.Println("too many shellgeis")
					return
				}
				result, err := DoShellGei(config, text)
				if err != nil {
					log.Println(err)
					return
				}
				if len(result) == 0 {
					return
				}
				err = TweetResult(api, tweet, result)
				if err != nil {
					log.Println(err)
					return
				}
			}()
		}
	}
}
