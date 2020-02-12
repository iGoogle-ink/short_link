package main

import (
	"log"
	"os"
)

type Env struct {
	S Storage
}

func getEnv() *Env {
	addr := os.Getenv("APP_REDIS_ADDR")
	if addr == "" {
		addr = "101.132.174.14:7001"
	}
	passwd := os.Getenv("APP_REDIS_PASSWD")
	log.Printf("connect redis (addr:%s, password:%s)\n", addr, passwd)
	redisCli := NewRedis(addr, passwd)
	return &Env{S: redisCli}
}
