// Copyright (C) edocseal. 2026-present.
//
// Created at 2026-02-07, by liasica

package snca

import (
	"testing"

	"go.uber.org/zap"
)

func TestRequest(t *testing.T) {
	failover := NewUrlFailover("http://localhost:5555", "https://www.baidu.com")

	l, _ := zap.NewDevelopment()
	zap.ReplaceGlobals(l)

	resp, err := createRestyClient(failover).R().Get("s?wd=“冰雪大国”是怎么炼成的&sa=fyb_n_homepage&rsv_dl=fyb_n_homepage&from=super&cl=3&tn=baidutop10&fr=top1000&rsv_idx=2&hisfilter=1")
	t.Log(resp, err)
}
