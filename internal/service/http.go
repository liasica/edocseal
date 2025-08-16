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

	"auroraride.com/edocseal"
	"auroraride.com/edocseal/internal/biz"
	"auroraride.com/edocseal/internal/g"
	"auroraride.com/edocseal/internal/task"
)

func StartHttpServer() {
	r := mux.NewRouter()
	r.HandleFunc("/s/{id}", shortUrl)
	r.HandleFunc("/maintain/stop/9geUbBHvX3caRWl1", stopTasks)
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

func getEnterpriseCert(w http.ResponseWriter, r *http.Request) {
	cfg := g.GetEnterpriseConfig()
	if cfg.Token != mux.Vars(r)["token"] {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// 如果名称相同返回200
	name := r.URL.Query().Get("name")
	serial := r.URL.Query().Get("serial")

	// 加载证书
	cert := cfg.GetCertificate()

	// 如果名称和序列号都相同，直接返回200
	if name == cert.Subject.CommonName && serial == cert.SerialNumber.String() {
		w.WriteHeader(http.StatusOK)
		return
	}

	// 如果名称不同，范围404
	if name != "" && name != cert.Subject.CommonName {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	ec := &edocseal.EnterpriseCertificate{}

	// 如果名称不同返回文件内容
	ec.Cert = string(cfg.GetCertificateBytes())

	b, err := os.ReadFile(cfg.PrivateKey)
	if err != nil {
		zap.L().Error("读取企业私钥失败", zap.Error(err), zap.String("private_key", cfg.PrivateKey))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	ec.Key = string(b)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	b, _ = jsoniter.Marshal(ec)
	_, _ = w.Write(b)
}
