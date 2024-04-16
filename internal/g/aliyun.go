// Copyright (C) edocseal. 2024-present.
//
// Created at 2024-04-16, by liasica

package g

import "strings"

func (a *AliyunOss) GetUrl(path string) (url string) {
	url = a.Url
	if !strings.HasSuffix(url, "/") {
		url += "/"
	}
	path = strings.TrimPrefix(path, "/")
	url += path
	return
}
