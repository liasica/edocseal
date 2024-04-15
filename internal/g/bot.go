// Copyright (C) edocseal. 2024-present.
//
// Created at 2024-04-15, by liasica

package g

import (
	"fmt"
	"os"

	bolt "go.etcd.io/bbolt"
)

var (
	boltDB         *bolt.DB
	ShortUrlBucket = []byte("SHORT_URL")
)

func NewBolt() *bolt.DB {
	if boltDB == nil {
		// 打开数据库
		var err error
		boltDB, err = bolt.Open(GetBboltPath(), 0600, nil)
		if err != nil {
			fmt.Printf("打开数据库失败: %v", err)
			os.Exit(1)
		}

		err = boltDB.Update(func(tx *bolt.Tx) error {
			if b := tx.Bucket(ShortUrlBucket); b == nil {
				_, err = tx.CreateBucket(ShortUrlBucket)
				return err
			}
			return nil
		})

		if err != nil {
			fmt.Printf("%s 创建失败: %v", ShortUrlBucket, err)
			os.Exit(1)
		}
	}

	return boltDB
}
