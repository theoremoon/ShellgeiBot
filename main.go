package main

import (
	"strconv"
	"bytes"
	"crypto/rand"
	"github.com/ChimeraCoder/anaconda"
	"golang.org/x/net/context"
	"log"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
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

	cmd := exec.Command("docker", "run", "--name", scriptFileName, "-v", scriptFilePath+":/"+scriptFileName, c.Image, "bash", "/"+scriptFileName)
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

func TweetResult(api *anaconda.TwitterApi, tweet anaconda.Tweet, result string) error {
	/// post done message
	v := url.Values{}
	_, err := api.PostTweet(result+" "+TweetUrl(tweet), v)
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

func IsMyTweet(api *anaconda.TwitterApi, tweet anaconda.Tweet) bool {
	v := url.Values{}
	self, err := api.GetSelf(v)
	if err != nil {
		return false
	}

	return self.Id == tweet.User.Id
}

func main() {
	if len(os.Args) < 7 {
		log.Println("6 arguments required. Consumer key, Consumer secret, Access token, Access token secret, timeout[sec], docker image name")
		return
	}
	anaconda.SetConsumerKey(os.Args[1])
	anaconda.SetConsumerSecret(os.Args[2])

	api := anaconda.NewTwitterApi(os.Args[3], os.Args[4])
	timeOutSec, err := strconv.Atoi(os.Args[5])
	if err != nil {
		return
	}
	image := os.Args[6]

	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	config := Config{
		image,
		dir,
		time.Duration(timeOutSec)*time.Second,
	}

	v := url.Values{}
	stream := api.UserStream(v)

	for {
		t := <-stream.C
		switch tweet := t.(type) {
		case anaconda.Tweet:
			go func() {
				is, text := IsShellGeiTweet(tweet.Text)
				if !is {
					return
				}
				is = IsMyTweet(api, tweet)
				if is {
					return
				}
				result, err := DoShellGei(config, text)
				if err != nil {
					return
				}
				if len(result) == 0 {
					return
				}
				TweetResult(api, tweet, result)
			}()
		}
	}
}
