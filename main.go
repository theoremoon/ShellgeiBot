// +build !test

package main

import (
  "database/sql"
  "github.com/ChimeraCoder/anaconda"
  _ "github.com/mattn/go-sqlite3"
  "net/url"
  "log"
  "os"
)

func ProcessTweet(tweet anaconda.Tweet, self anaconda.User, api *anaconda.TwitterApi, db *sql.DB, botConfig BotConfig) {

  // check valid shellgei tweet
  if tweet.RetweetedStatus != nil {
    return
  }
  is := IsShellGeiTweet(tweet.FullText, botConfig.Tags)
  if !is {
    return
  }
  if self.Id == tweet.User.Id {
    return
  }
  if !IsFollower(api, tweet) {
    return
  }

  text := ExtractShellgei(tweet, self, api, botConfig.Tags)
  t, err := tweet.CreatedAtTime()
  if err != nil{
    log.Println(err)
    return
  }

  InsertShellGei(db, tweet.User.Id, tweet.User.ScreenName, tweet.Id, text, t.Unix())

  result, err := RunCmd(text, botConfig)
  InsertResult(db, tweet.Id, result, err)

  if err != nil {
    if err.(*StdError) == nil {
      _, _ = api.PostTweet("@theoldmoon0602 internal error", url.Values{})
    }
    return
  }

  if len(result) == 0 {
    return
  }

  err = TweetResult(api, tweet, MakeTweetable(result))
  if err != nil {
    log.Println(err)
  }
  return
}

func main() {
  if len(os.Args) < 3 {
    log.Fatalf("<Usage>%s: TwitterConfig.json ShellgeiConfig.json", os.Args[0])
  }

  twitterKey, err := ParseTwitterKey(os.Args[1])
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

  stream := api.UserStream(v)
  for {
    t := <-stream.C
    switch tweet := t.(type) {
    case anaconda.Tweet:
      botConfig, err := ParseBotConfig(os.Args[2])
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
