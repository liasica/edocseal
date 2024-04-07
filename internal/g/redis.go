// Copyright (C) edocseal. 2024-present.
//
// Created at 2024-04-07, by liasica

package g

import "github.com/redis/go-redis/v9"

var rdb *redis.Client

// NewRedis 初始化redis
func NewRedis() *redis.Client {
	if rdb != nil {
		return rdb
	}
	rdb = redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Addr,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})
	return rdb
}
