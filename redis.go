package main

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis/v7"
)

const (
	// URLIDKEY is global counter
	URLIDKEY = "next.url.id"

	// ShortlinkKey =
	ShortlinkKey = "shortlink:%s:url"

	// URLHashKey
	URLHashKey = "URLHashKey:%s:url"

	// ShortlinkDetailKey
	ShortlinkDetailKey = "ShortlinkDetailKey:%s:detail"
)

type RedisClient struct {
	Cli *redis.ClusterClient
}

type URLDetail struct {
	URL           string        `json:"url"`
	CreatedAt     string        `json:"created_at"`
	ExpireMinutes time.Duration `json:"expire_minutes"`
}

func NewRedis(addr string, passwd string) *RedisClient {
	cli := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:    []string{addr},
		Password: passwd,
	})
	ping, err := cli.Ping().Result()
	if err != nil {
		panic(err)
	}
	log.Println("Redis PING:", ping)
	return &RedisClient{Cli: cli}
}

func (rc *RedisClient) Shorten(url string, expire int64) (string, error) {
	h := getHash(url)
	short, err := rc.Cli.Get(fmt.Sprintf(URLHashKey, h)).Result()
	if err != nil {
		if err != redis.Nil {
			return "", err
		}
	}
	if short != "" && short != "{}" {
		return short, nil
	}

	err = rc.Cli.Incr(URLIDKEY).Err()
	if err != nil {
		return "", err
	}

	id, err := rc.Cli.Get(URLIDKEY).Result()
	if err != nil {
		return "", err
	}
	log.Println("id:", id)

	short = base64.StdEncoding.EncodeToString([]byte(id))
	err = rc.Cli.Set(fmt.Sprintf(ShortlinkKey, short), url, time.Minute*time.Duration(expire)).Err()
	if err != nil {
		return "", err
	}

	err = rc.Cli.Set(fmt.Sprintf(URLHashKey, h), short, time.Minute*time.Duration(expire)).Err()
	if err != nil {
		return "", err
	}

	detail := &URLDetail{
		URL:           url,
		CreatedAt:     time.Now().String(),
		ExpireMinutes: time.Duration(expire),
	}
	bs, err := json.Marshal(detail)
	if err != nil {
		return "", err
	}

	err = rc.Cli.Set(fmt.Sprintf(ShortlinkDetailKey, short), bs, time.Minute*time.Duration(expire)).Err()
	if err != nil {
		return "", err
	}
	return short, nil
}

func (rc *RedisClient) ShortlinkInfo(short string) (interface{}, error) {
	shortDetail, err := rc.Cli.Get(fmt.Sprintf(ShortlinkDetailKey, short)).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, &StatusError{Code: 404, Err: errors.New("Unknown Short URL.")}
		}
		return nil, err
	}
	return shortDetail, nil
}

func (rc *RedisClient) Unshorten(short string) (string, error) {
	url, err := rc.Cli.Get(fmt.Sprintf(ShortlinkKey, short)).Result()
	if err != nil {
		if err == redis.Nil {
			return "", &StatusError{Code: 404, Err: errors.New("Unknown Short URL.")}
		}
		return "", err
	}
	return url, nil
}

func getHash(url string) string {
	hash := md5.New()
	hash.Write([]byte(url))
	return hex.EncodeToString(hash.Sum(nil))
}
