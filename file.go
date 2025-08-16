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
	"strings"

	"go.uber.org/zap/buffer"
)

// CreateDirectory 文件夹创建
func CreateDirectory(dir string) (err error) {
	_, err = os.Stat(dir)
	if os.IsNotExist(err) {
		return os.MkdirAll(dir, os.ModePerm)
	}
	return
}

// FileExists 文件是否存在
func FileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// FileNameWithoutExtension 获取文件名
func FileNameWithoutExtension(fileName string) string {
	fileName = filepath.Base(fileName)
	if pos := strings.LastIndexByte(fileName, '.'); pos != -1 {
		return fileName[:pos]
	}
	return fileName
}

// DownloadFile 下载文件
func DownloadFile(url, dir string) (sn, path string, err error) {
	var resp *http.Response
	resp, err = http.Get(url)
	if err != nil {
		return
	}

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("文档请求失败，StatusCode = %s", resp.Status)
		return
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
	path = dir + "/" + sn + ext
	err = os.WriteFile(path, buf.Bytes(), os.ModePerm)
	return
}

// FileMd5 文件md5
func FileMd5(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	hash := md5.New()
	if _, err = io.Copy(hash, file); err != nil {
		return "", err
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}

// FileCopy 文件复制
func FileCopy(src, dst string) error {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return fmt.Errorf("%s is not a regular file", src)
	}

	var (
		source      *os.File
		destination *os.File
	)

	source, err = os.Open(src)
	if err != nil {
		return err
	}
	defer func(source *os.File) {
		_ = source.Close()
	}(source)

	destination, err = os.Create(dst)
	if err != nil {
		return err
	}
	defer func(destination *os.File) {
		_ = destination.Close()
	}(destination)
	_, err = io.Copy(destination, source)

	return err
}

// DeleteDirectory 删除文件夹
func DeleteDirectory(dir string) (err error) {
	return os.RemoveAll(dir)
}

// GetFileName 获取文件名
// 如果params不为空，则使用params[0]作为扩展名
func GetFileName(filename string, params ...string) (name string) {
	basename := filepath.Base(filename)
	ext := filepath.Ext(basename)
	if len(params) > 0 {
		ext = params[0] + ext
	}
	name = strings.TrimSuffix(basename, ext)
	return
}
