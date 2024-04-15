// Copyright (C) edocseal. 2024-present.
//
// Created at 2024-04-15, by liasica

package g

import (
	"github.com/teris-io/shortid"
)

var sid *shortid.Shortid

func NewShortId() *shortid.Shortid {
	if sid == nil {
		sid = shortid.MustNew(1, shortid.DefaultABC, 2342)
	}
	return sid
}
