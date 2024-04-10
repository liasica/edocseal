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

type AliyunOss struct {
	AccessKeyId     string
	AccessKeySecret string
	Bucket          string
	Endpoint        string
	Url             string
}

type Config struct {
	// 企业签章图片
	Seal string

	// 目录配置
	Dir struct {
		Template string // 模板目录
		Runtime  string // 运行时目录
		Document string // 文档目录
	}

	// 文档配置
	Document struct {
		BucketPath string // OSS存储路径
	}

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

	Aliyun struct {
		Oss AliyunOss
	}
}

func readConfig() (err error) {
	err = viper.ReadInConfig()
	if err != nil {
		fmt.Printf("配置读取失败: %s\n", err)
		return
	}

	cfg = &Config{}
	err = viper.Unmarshal(cfg)
	if err != nil {
		return
	}

	// 获取签章完整路径
	cfg.Seal, err = filepath.Abs(cfg.Seal)

	return
}

// LoadConfig 加载配置文件
func LoadConfig(configFile string) {
	path = filepath.Dir(configFile)

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

	// 监听配置文件变动后重载
	viper.OnConfigChange(func(in fsnotify.Event) {
		fmt.Println("配置文件变动，重新加载")
		_ = readConfig()
	})

	viper.WatchConfig()
}

func GetConfig() *Config {
	return cfg
}

// GetSeal 获取企业签章图片
func GetSeal() string {
	return cfg.Seal
}

// GetConfigPath 获取配置目录
func GetConfigPath() string {
	return path
}

// GetRPCBind 获取RPC绑定地址
func GetRPCBind() string {
	return cfg.RPC.Bind
}

// GetTemplateDir 获取模板路径
func GetTemplateDir() string {
	return cfg.Dir.Template
}

// GetRuntimeDir 获取运行时目录
func GetRuntimeDir() string {
	return cfg.Dir.Runtime
}

// GetDocumentDir 获取文档目录
func GetDocumentDir() string {
	return cfg.Dir.Document
}

// GetAliyunOss 获取阿里云OSS配置
func GetAliyunOss() AliyunOss {
	return cfg.Aliyun.Oss
}
