// Copyright (C) edocseal. 2024-present.
//
// Created at 2024-04-13, by liasica

package service

import (
	"fmt"
	"net/http"
	"os"
	"sync"

	"github.com/gorilla/mux"
	"go.uber.org/zap"

	"auroraride.com/edocseal/internal/biz"
	"auroraride.com/edocseal/internal/g"
	"auroraride.com/edocseal/internal/task"
)

func StartHttpServer() {
	r := mux.NewRouter()
	r.HandleFunc("/s/{id}", shortUrl)
	r.HandleFunc("/maintain/stop/9geUbBHvX3caRWl1", stopTasks)
	zap.L().Info("API启动", zap.String("bind", g.GetHttpBind()))
	err := http.ListenAndServe(g.GetHttpBind(), r)
	if err != nil {
		fmt.Printf("HTTP服务启动失败：%s\n", err)
		os.Exit(1)
	}
}

// 停止所有任务
func stopTasks(w http.ResponseWriter, _ *http.Request) {
	wg := &sync.WaitGroup{}
	wg.Add(2)

	go task.DocumentTask().HaltAndWait(wg)
	go task.SignTask().HaltAndWait(wg)

	wg.Wait()
	_, _ = w.Write([]byte("ok\n"))
}

func shortUrl(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		http.NotFound(w, r)
		return
	}
	src, err := biz.GetShortUrl(id)
	if err != nil {
		zap.L().Error("获取短链接失败", zap.Error(err), zap.String("id", id))
	}
	// TODO: oss保存到私有bucket并且使用acl生成临时访问链接后跳转
	http.Redirect(w, r, src, http.StatusMovedPermanently)
}
