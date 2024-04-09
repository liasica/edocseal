// Copyright (C) edocseal. 2024-present.
//
// Created at 2024-04-08, by liasica

package edocseal

import (
	"strconv"
)

const (
	ErrDefaultCode = iota + 11000
	ErrNotFoundCode
	ErrInvalidArgumentCode
)

type Error struct {
	Code    int
	Message string
}

func (e *Error) Error() string {
	return strconv.Itoa(e.Code) + ": " + e.Message
}

func NewError(code int, message string) *Error {
	return &Error{Code: code, Message: message}
}

func ErrorWrapping(err error) *Error {
	return NewError(ErrDefaultCode, err.Error())
}

var (
	ErrFieldNotFound = func(field string) *Error {
		return NewError(ErrNotFoundCode, field+" 字段不存在")
	}
)
