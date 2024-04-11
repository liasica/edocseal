// Copyright (C) edocseal. 2024-present.
//
// Created at 2024-04-11, by liasica

package g

import (
	"strconv"

	"github.com/sony/sonyflake"
)

var st = sonyflake.NewSonyflake(sonyflake.Settings{})

func GetSnowflake() *sonyflake.Sonyflake {
	return st
}

func GetID() string {
	id, _ := st.NextID()
	return strconv.FormatUint(id, 10)
}
