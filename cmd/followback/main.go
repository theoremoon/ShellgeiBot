/// N日以上フォロワーだったひとをフォローバックする仕組み
/// フォロー解除されていたら解除し返す
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	twitterkey "github.com/theoldmoon0602/ShellgeiBot/twitter"
	"golang.org/x/xerrors"
)

type Follower struct {
	UserID       int64 `json:"userid"`
	FollowedFrom int64 `json:"followedfrom"`
}

func makeClient(key twitterkey.TwitterKey) *twitter.Client {
	config := oauth1.NewConfig(key.ConsumerKey, key.ConsumerSecret)
	token := oauth1.NewToken(key.AccessToken, key.AccessSecret)
	httpClient := config.Client(oauth1.NoContext, token)

	return twitter.NewClient(httpClient)
}

func loadFollowers(filename string) ([]Follower, error) {
	if _, err := os.Stat(filename); err != nil {
		return []Follower{}, nil
	}

	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, xerrors.Errorf(": %w", err)
	}
	followers := make([]Follower, 0)
	err = json.Unmarshal(data, &followers)
	if err != nil {
		return nil, xerrors.Errorf(": %w", err)
	}
	return followers, nil
}

func listCurrentFollowerIDs(client *twitter.Client) ([]int64, error) {
	var cursor int64 = -1
	ids := make([]int64, 0)
	for cursor != 0 {
		followerChunk, _, err := client.Followers.IDs(&twitter.FollowerIDParams{
			Cursor: cursor,
			Count:  5000, // twitter APIで指定できる最大値
		})
		if err != nil {
			return nil, xerrors.Errorf(": %w", err)
		}
		ids = append(ids, followerChunk.IDs...)
		cursor = followerChunk.NextCursor
	}
	return ids, nil
}

func listCurrentFollowingUserIDs(client *twitter.Client) ([]int64, error) {
	var cursor int64 = -1
	ids := make([]int64, 0)
	for cursor != 0 {
		idChunk, _, err := client.Friends.IDs(&twitter.FriendIDParams{
			Cursor: cursor,
			Count:  5000, // twitter APIで指定できる最大値
		})
		if err != nil {
			return nil, xerrors.Errorf(": %w", err)
		}
		ids = append(ids, idChunk.IDs...)
		cursor = idChunk.NextCursor
	}
	return ids, nil
}

/// xsに存在するけどysに存在しないやつの集合を求める
func difference(xs, ys []int64) []int64 {
	set := make(map[int64]struct{})
	for _, y := range ys {
		set[y] = struct{}{}
	}

	zs := make([]int64, 0, len(xs))
	for _, x := range xs {
		if _, e := set[x]; !e {
			zs = append(zs, x)
		}
	}
	return zs
}

/// idsに渡されたidを持つuserをunfollowする
func unfollowByIDs(client *twitter.Client, ids []int64) error {
	// twitter APIの制限により400件より多くを一度に処理しない
	// usage的に黙って400件だけ処理しても困らないのでassertとか入れてない
	n := len(ids)
	if n > 400 {
		n = 400
	}
	for _, id := range ids[:n] {
		_, _, err := client.Friendships.Destroy(&twitter.FriendshipDestroyParams{
			UserID: id,
		})
		if err != nil {
			return xerrors.Errorf(": %w", err)
		}
	}
	return nil
}

/// idsに渡されたidを持つuserをfollowする
func followByIDs(client *twitter.Client, ids []int64) error {
	// twitter APIの制限により400件より多くを一度に処理しない
	// usage的に黙って400件だけ処理しても困らないのでassertとか入れてない
	n := len(ids)
	if n > 400 {
		n = 400
	}
	for _, id := range ids[:n] {
		_, _, err := client.Friendships.Create(&twitter.FriendshipCreateParams{
			UserID: id,
		})
		if err != nil {
			return xerrors.Errorf(": %w", err)
		}
	}
	return nil
}

func run() error {
	var followerFile string
	var keyFile string
	var ndays int
	flag.StringVar(&followerFile, "followers", "", "<followers.json>")
	flag.StringVar(&keyFile, "twitter", "", "<twitterkey.json>")
	flag.IntVar(&ndays, "ndays", 1, "<days to followback>")
	flag.Usage = func() {
		fmt.Printf("Usage: %s\n\n", os.Args[1])
		flag.PrintDefaults()
	}
	flag.Parse()
	if ndays <= 0 {
		flag.Usage()
		return nil
	}

	// load past followers
	pastFollowers, err := loadFollowers(followerFile)
	if err != nil {
		return xerrors.Errorf(": %w", err)
	}
	pastFollowerMap := make(map[int64]Follower)
	for _, u := range pastFollowers {
		pastFollowerMap[u.UserID] = u
	}

	// initialize twitter client
	key, err := twitterkey.ParseTwitterKey(keyFile)
	if err != nil {
		return xerrors.Errorf(": %w", err)
	}
	client := makeClient(key)

	// get current followers
	followerIDs, err := listCurrentFollowerIDs(client)
	if err != nil {
		return xerrors.Errorf(": %w", err)
	}

	// get current followings
	followingIDs, err := listCurrentFollowingUserIDs(client)
	if err != nil {
		return xerrors.Errorf(": %w", err)
	}

	// フォローするユーザ決めたりする
	now := time.Now()
	t := now.AddDate(0, 0, -ndays).Unix() // この時刻より前にフォローしてくれている必要がある
	newFollowerIDs := make([]int64, 0)
	newFollowingIDs := make([]int64, 0)
	for _, id := range followerIDs {
		u, e := pastFollowerMap[id]
		if e && u.FollowedFrom < t && len(newFollowingIDs) < 400 {
			newFollowingIDs = append(newFollowingIDs, id)
			delete(pastFollowerMap, id)
		} else if !e {
			newFollowerIDs = append(newFollowerIDs, id)
		}
	}

	// 新しくフォローされていたユーザの情報を保存しておく
	saveFollowers := make([]Follower, 0)
	for _, u := range pastFollowerMap {
		saveFollowers = append(saveFollowers, u)
	}
	for _, id := range newFollowerIDs {
		saveFollowers = append(saveFollowers, Follower{
			UserID:       id,
			FollowedFrom: now.Unix(),
		})
	}
	jsonBytes, err := json.MarshalIndent(saveFollowers, "", "  ")
	if err != nil {
		return xerrors.Errorf(": %w", err)
	}
	if err := os.WriteFile(followerFile, jsonBytes, 0755); err != nil {
		return xerrors.Errorf(": %w", err)
	}

	wg := &sync.WaitGroup{}
	// unfollow
	wg.Add(1)
	go func() {
		unfollowIDs := difference(followingIDs, followerIDs)
		err := unfollowByIDs(client, unfollowIDs)
		if err != nil {
			log.Printf("%+v\n", err)
		}
		wg.Done()
	}()

	// follow
	wg.Add(1)
	go func() {
		err := followByIDs(client, newFollowingIDs)
		if err != nil {
			log.Printf("%+v\n", err)
		}
		wg.Done()
	}()

	wg.Wait()
	return nil
}

func main() {
	if err := run(); err != nil {
		log.Fatalf("%+v\n", err)
	}
}
