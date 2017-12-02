package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"
)

type BotConfigJson struct {
	DockerImage string `json:"dockerimage"`
	Workdir     string `json:"workdir"`
	Timeout     string `json:"timeout"`
}

type BotConfig struct {
	DockerImage string
	Workdir     string
	Timeout     time.Duration
}

type TwitterKeys struct {
	ConsumerKey    string `json:"ConsumerKey"`
	ConsumerSecret string `json:"ConsumerSecret"`
	AccessToken    string `json:"AccessToken"`
	AccessSecret   string `json:"AccessSecret"`
}

func ParseBotConfig(file string) (BotConfig, error) {
	var c BotConfigJson
	var config BotConfig

	// read json
	raw, err := ioutil.ReadFile(file)
	if err != nil {
		return config, err
	}
	err = json.Unmarshal(raw, &c)
	if err != nil {
		return config, err
	}

	// convert json to config type
	config.DockerImage = c.DockerImage
	config.Workdir, err = filepath.Abs(c.Workdir)
	if err != nil {
		return config, err
	}
	config.Timeout, err = time.ParseDuration(c.Timeout)
	if err != nil {
		return config, err
	}
	return config, nil
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

func main() {
	if len(os.Args) < 3 {
		log.Fatalf("<Usage>%s: TwitterConfig.json ShellgeiConfig.json", os.Args[0])
	}

	twitterKey, err := ParseTwitterKey(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	botConfig, err := ParseBotConfig(os.Args[2])
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(twitterKey.ConsumerSecret)
	fmt.Println(botConfig.Timeout)
	fmt.Println(botConfig.Workdir)
}
