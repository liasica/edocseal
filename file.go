// Copyright (C) edocseal. 2024-present.
//
// Created at 2024-04-08, by liasica

package edocseal

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"go.uber.org/zap/buffer"
)

// DownloadFile 下载文件
func DownloadFile(url, dir string) (err error, sn string) {
	// 判定文件夹是否存在，不存在则创建
	_, err = os.Stat(dir)
	if os.IsNotExist(err) {
		err = os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			return
		}
	}

	var resp *http.Response
	resp, err = http.Get(url)
	if err != nil {
		return
	}

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("文档请求失败，StatusCode = %s", resp.Status), ""
	}

	// 获取文件md5
	buf := new(buffer.Buffer)
	hash := md5.New()
	tee := io.TeeReader(resp.Body, hash)
	_, err = io.Copy(buf, tee)
	if err != nil {
		return
	}

	sn = hex.EncodeToString(hash.Sum(nil))
	ext := filepath.Ext(url)

	// 存储文件
	err = os.WriteFile(dir+"/"+sn+ext, buf.Bytes(), os.ModePerm)
	return
}
