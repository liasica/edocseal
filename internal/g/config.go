// Copyright (C) edocseal. 2024-present.
//
// Created at 2024-04-05, by liasica

package g

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"

	"github.com/liasica/edocseal"
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

type Snca struct {
	Url string
}

type Config struct {
	// Signer路径
	Signer string

	// 是否自签名
	SelfSign bool

	// 同时处理任务数
	Task struct {
		Sign     int // 签约任务
		Document int // 文档任务
	}

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
	Enterprise struct {
		Seal        string // 企业签章图片
		Certificate string // 企业证书
		PrivateKey  string // 企业私钥
	}

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

	Snca Snca
}

func readConfig() (err error) {
	err = viper.ReadInConfig()
	if err != nil {
		return fmt.Errorf("配置读取失败: %s\n", err)
	}

	cfg = &Config{}
	err = viper.Unmarshal(cfg)
	if err != nil {
		return
	}

	// 获取企业签名文件完整路径
	cfg.Enterprise.Seal, _ = filepath.Abs(cfg.Enterprise.Seal)
	if !edocseal.FileExists(cfg.Enterprise.Seal) {
		return errors.New("企业签章图片不存在")
	}
	cfg.Enterprise.PrivateKey, _ = filepath.Abs(cfg.Enterprise.PrivateKey)
	if !edocseal.FileExists(cfg.Enterprise.PrivateKey) {
		return errors.New("企业私钥不存在")
	}
	cfg.Enterprise.Certificate, _ = filepath.Abs(cfg.Enterprise.Certificate)
	if !edocseal.FileExists(cfg.Enterprise.Certificate) {
		return errors.New("企业证书不存在")
	}
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
		fmt.Printf("配置读取失败: %v\n", err)
		os.Exit(1)
	}

	if cfg.Task.Sign == 0 {
		cfg.Task.Sign = 3
	}
	if cfg.Task.Document == 0 {
		cfg.Task.Document = 3
	}
}

// GetSignTaskNum 获取签约任务数
func GetSignTaskNum() int {
	return cfg.Task.Sign
}

// GetDocumentTaskNum 获取文档任务数
func GetDocumentTaskNum() int {
	return cfg.Task.Document
}

// GetSeal 获取企业签章图片
func GetSeal() string {
	return cfg.Enterprise.Seal
}

// GetCertificate 获取企业证书
func GetCertificate() string {
	return cfg.Enterprise.Certificate
}

// GetPrivateKey 获取企业私钥
func GetPrivateKey() string {
	return cfg.Enterprise.PrivateKey
}

// GetSigner 获取Signer路径
func GetSigner() string {
	return cfg.Signer
}

// IsSelfSign 是否自签名
func IsSelfSign() bool {
	return cfg.SelfSign
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

// GetSnca 获取SNCA配置
func GetSnca() *Snca {
	return &cfg.Snca
}
