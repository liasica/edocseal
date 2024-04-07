// Copyright (C) edocseal. 2024-present.
//
// Created at 2024-04-05, by liasica

package g

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var (
	cfg  *Config
	path string
)

type CertificatePath struct {
	Certificate string
	PrivateKey  string
}

type Config struct {
	// 根证书和私钥，用于签发证书
	RootCertificate CertificatePath

	// 企业证书，用于签署协议
	EnterpriseCertificate CertificatePath

	// 日志配置
	Logger struct {
		Console  bool   // 是否输出至控制台
		Redis    bool   // 是否输出至Redis
		RedisKey string // Redis输出key
	}

	// RPC配置
	RPC struct {
		Bind string // 绑定地址
	}

	// Redis配置
	Redis struct {
		Addr     string // 地址
		Password string // 密码
		DB       int    // 数据库
	}
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

	path = filepath.Dir(configFile)

	// 监听配置文件变动后重载
	viper.OnConfigChange(func(in fsnotify.Event) {
		fmt.Println("配置文件变动，重新加载")
		_ = readConfig()
	})

	viper.WatchConfig()
}

// GetConfigPath 获取配置目录
func GetConfigPath() string {
	return path
}

// GetRPCBind 获取RPC绑定地址
func GetRPCBind() string {
	return cfg.RPC.Bind
}
