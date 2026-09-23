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
	jsoniter "github.com/json-iterator/go"
	"go.uber.org/zap"

	"auroraride.com/edocseal/internal/biz"
	"auroraride.com/edocseal/internal/g"
	"auroraride.com/edocseal/internal/task"
)

// maintainToken 维护接口路径中的令牌
const maintainToken = "9geUbBHvX3caRWl1"

func StartHttpServer() {
	r := mux.NewRouter()
	r.HandleFunc("/s/{id}", shortUrl)
	r.HandleFunc("/maintain/stop/"+maintainToken, stopTasks)
	r.HandleFunc("/enterprise/cert/{token}", getEnterpriseCert).Methods("GET")
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

// 按统一社会信用代码返回企业陕西CA 证书与私钥，供其他环境拉取
func getEnterpriseCert(w http.ResponseWriter, r *http.Request) {
	token := g.GetEnterpriseConfig().Token
	if token == "" || token != mux.Vars(r)["token"] {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	code := r.URL.Query().Get("code")
	result, err := biz.EnterpriseCertificatePEM(code)
	if err != nil {
		zap.L().Error("获取企业证书失败", zap.Error(err), zap.String("code", code))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	b, _ := jsoniter.Marshal(result)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(b)
}
