//go:build !test
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
	"sync"
	"time"

	"github.com/ChimeraCoder/anaconda"
	"github.com/mattn/go-sixel"
	_ "github.com/mattn/go-sqlite3"
	"github.com/theoldmoon0602/ShellgeiBot/twitter"
)

func processTweet(tweet anaconda.Tweet, self anaconda.User, api *anaconda.TwitterApi, db *sql.DB, config botConfig) {
	// check if it is valid shellgei tweet
	if tweet.RetweetedStatus != nil {
		return
	}
	if !isShellGeiTweet(tweet, config.Tags) {
		return
	}
	log.Printf("treat %d as the shellgei tweet \n", tweet.Id)
	if self.Id == tweet.User.Id {
		return
	}
	if !isFollower(api, tweet) {
		return
	}
	if isProcessed(db, tweet.Id) {
		log.Printf("%d: already processed\n", tweet.Id)
		return
	}

	t, err := tweet.CreatedAtTime()
	if err != nil {
		log.Println(err)
		return
	}
	text, mediaUrls, err := extractShellgei(tweet, self, api, config.Tags, []int64{})
	if err != nil {
		log.Println(err)
		return
	}

	insertShellGei(db, tweet.User.Id, tweet.User.ScreenName, tweet.Id, text, t.Unix())

	result, b64imgs, err := runCmd(text, mediaUrls, config)
	result = makeTweetable(result)
	insertResult(db, tweet.Id, result, err)

	if err != nil {
		log.Printf("runcmd: %v\n", err)
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

// / ShellgeiBot main function
func botMain(twitterConfigFile, botConfigFile string) {
	twitterKey, err := twitter.ParseTwitterKey(twitterConfigFile)
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
	api.SetLogger(anaconda.BasicLogger)

	self, err := api.GetSelf(url.Values{})
	if err != nil {
		log.Fatal(err)
	}

	if err != nil {
		log.Fatal(err)
	}

	config, err := parseBotConfig(botConfigFile)
	ticker := time.NewTicker(1 * time.Minute)
	v := url.Values{}
	v.Set("count", "200")
	for {
		tweets, err := api.GetHomeTimeline(v)
		if err != nil {
			log.Printf("GetHomeTimeline: %v", err)
			<-ticker.C
			continue
		}

		for _, tweet := range tweets {
			processTweet(tweet, self, api, db, config)
			v.Set("since_id", tweet.IdStr)
		}

		<-ticker.C
	}
}

func botTest(botConfigFile string, scripts []string) {
	config, err := parseBotConfig(botConfigFile)
	if err != nil {
		log.Fatal(err)
	}

	type workResult struct {
		Stdout string
		Images []image.Image
		Time   time.Duration
		Error  error
	}

	worker := func(scriptFile string, result chan<- workResult, wg *sync.WaitGroup) {
		defer wg.Done()

		script, err := ioutil.ReadFile(scriptFile)
		if err != nil {
			result <- workResult{
				Error: err,
			}
			return
		}

		start := time.Now()
		stdout, b64imgs, err := runCmd(string(script), []string{}, config)
		t := time.Since(start)
		if err != nil {
			result <- workResult{
				Error: err,
				Time:  t,
			}
			return
		}
		stdout = makeTweetable(stdout)
		if stdout == "" && len(b64imgs) == 0 {
			result <- workResult{
				Error: fmt.Errorf("Empty result"),
				Time:  t,
			}
			return
		}

		images := make([]image.Image, 0, 4)
		for _, b64img := range b64imgs {
			imgBytes, err := base64.StdEncoding.DecodeString(b64img)
			if err != nil {
				// if media is not a valid image (e.g. meaningless bytestream, video)
				log.Println(err)
				continue
			}

			img, _, err := image.Decode(bytes.NewReader(imgBytes))
			if err != nil {
				log.Println(err)
				continue
			}
			images = append(images, img)
		}

		result <- workResult{
			Stdout: stdout,
			Images: images,
			Time:   t,
			Error:  nil,
		}
	}

	var wg sync.WaitGroup
	results := make(chan workResult, len(scripts))
	for _, scriptFile := range scripts {
		wg.Add(1)
		go worker(scriptFile, results, &wg)
	}

	wg.Wait()

	for i := 1; i <= len(scripts); i++ {
		r := <-results

		fmt.Printf("Result: %d\n", i)
		fmt.Println("=== Stdout ===")
		fmt.Println(r.Stdout)
		fmt.Println("=== Images ===")
		for _, img := range r.Images {
			sixel.NewEncoder(os.Stdout).Encode(img)
		}
		fmt.Println("=== Error  ===")
		fmt.Println(r.Error)
		fmt.Println("===  Time  ===")
		fmt.Println(r.Time)
		fmt.Println()
	}
}

func main() {
	if len(os.Args) < 3 {
		log.Fatalf("<Usage>%s: TwitterConfig.json ShellgeiConfig.json | -test ShellgeiConfig.json [scripts]", os.Args[0])
	}

	if os.Args[1] == "-test" {
		// testing mode
		botTest(os.Args[2], os.Args[3:])
	} else {
		// normal mode
		botMain(os.Args[1], os.Args[2])
	}
}
