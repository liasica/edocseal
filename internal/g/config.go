// Copyright (C) edocseal. 2024-present.
//
// Created at 2024-04-05, by liasica

package g

import (
	"crypto/rsa"
	"crypto/x509"
	"fmt"
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var cfg *Config

type Config struct {
	// 根证书和私钥
	Root *Certificate

	// 配置目录
	Path string

	// 日志配置
	Logger struct {
		// 是否输出至控制台
		Console bool
		// 是否输出至Redis
		Redis bool
	}
}

// Certificate 证书配置
type Certificate struct {
	Certificate *x509.Certificate
	PrivateKey  *rsa.PrivateKey
}

func readConfig() (err error) {
	err = viper.ReadInConfig()
	if err != nil {
		fmt.Printf("配置读取失败: %s\n", err)
		return
	}

	cfg = &Config{}
	return viper.Unmarshal(cfg)
}

// LoadConfig 加载配置文件
func LoadConfig(configFile string) {
	// 判定配置文件是否存在
	_, err := os.Stat(configFile)
	if err != nil {
		fmt.Println("配置文件不存在")
		os.Exit(1)
	}

	viper.SetConfigFile(configFile)
	viper.AutomaticEnv()

	// 读取配置文件
	err = readConfig()
	if err != nil {
		os.Exit(1)
	}

	cfg.Path = filepath.Dir(configFile)

	// 监听配置文件变动后重载
	viper.OnConfigChange(func(in fsnotify.Event) {
		fmt.Println("配置文件变动，重新加载")
		_ = readConfig()
	})

	viper.WatchConfig()
}

// GetRootCertificate 获取根证书
func GetRootCertificate() *x509.Certificate {
	return cfg.Root.Certificate
}

// GetRootPrivateKey 获取根证书私钥
func GetRootPrivateKey() *rsa.PrivateKey {
	return cfg.Root.PrivateKey
}

// GetConfigPath 获取配置目录
func GetConfigPath() string {
	return cfg.Path
}
