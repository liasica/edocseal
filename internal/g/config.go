// Copyright (C) edocseal. 2024-present.
//
// Created at 2024-04-05, by liasica

package g

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/viper"

	"auroraride.com/edocseal"
)

var (
	cfg        *Config
	configFile string
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
	Url          string
	UrlFallback  string
	Source       string
	CustomerType string
}

type Enterprise struct {
	Seal        string // 企业签章图片
	Certificate string // 企业证书
	PrivateKey  string // 企业私钥
	PersonName  string // 代办姓名
	Phone       string // 代办电话
	Idcard      string // 代办身份证
	Province    string // 省份
	City        string // 城市
	CreditCode  string // 统一社会信用代码
	Name        string // 企业名称
}

type Config struct {
	// Signer路径
	Signer string

	// 是否自签名
	SelfSign bool

	// 短链接前缀
	ShortUrlPrefix string

	// Bbolt路径
	Bbolt struct {
		Path string
	}

	// 同时处理任务数
	Task struct {
		Sign     int // 签约任务
		Document int // 文档任务
	}

	// 目录配置
	Dir struct {
		Template    string // 模板目录
		Runtime     string // 运行时目录
		Document    string // 文档目录
		Certificate string // 证书目录
	}

	// 文档配置
	Document struct {
		BucketPath string // OSS存储路径
	}

	// 根证书和私钥，用于签发证书
	RootCertificate CertificatePath

	// 企业证书，用于签署协议
	Enterprise *Enterprise

	// 日志配置
	Logger struct {
		Console    bool   // 是否输出至控制台
		Redis      bool   // 是否输出至Redis
		RedisKey   string // Redis输出key
		LoggerName string // 日志名称
	}

	// RPC配置
	RPC struct {
		Bind string // 绑定地址
	}

	// HTTP配置
	Http struct {
		Bind string // 绑定地址
	}

	// Redis配置
	Redis struct {
		Addr     string // 地址
		Password string // 密码
		DB       int    // 数据库
	}

	Aliyun struct {
		Oss *AliyunOss
	}

	Snca Snca

	Postgres struct {
		Dsn   string
		Debug bool
	}
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
	if !edocseal.FileExists(cfg.Enterprise.Seal) {
		return errors.New("企业签章图片不存在")
	}
	if !edocseal.FileExists(cfg.Enterprise.PrivateKey) {
		return errors.New("企业私钥不存在")
	}
	if !edocseal.FileExists(cfg.Enterprise.Certificate) {
		return errors.New("企业证书不存在")
	}
	return
}

// LoadConfig 加载配置文件
func LoadConfig(path string) {
	configFile = path

	// 判定配置文件是否存在
	_, err := os.Stat(path)
	if err != nil {
		fmt.Println("配置文件不存在")
		os.Exit(1)
	}

	viper.SetConfigFile(path)
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

// GetEnterpriseConfig 获取企业配置
func GetEnterpriseConfig() *Enterprise {
	return cfg.Enterprise
}

// UpdateEnterpriseConfig 更新企业配置
func UpdateEnterpriseConfig(key, cert string) {
	cfg.Enterprise.PrivateKey = key
	cfg.Enterprise.Certificate = cert
}

// GetShortUrlPrefix 获取短链接前缀
func GetShortUrlPrefix() string {
	return cfg.ShortUrlPrefix
}

// GetBboltPath 获取Bbolt路径
func GetBboltPath() string {
	return cfg.Bbolt.Path
}

// GetSigner 获取Signer路径
func GetSigner() string {
	return cfg.Signer
}

// IsSelfSign 是否自签名
func IsSelfSign() bool {
	return cfg.SelfSign
}

// GetRPCBind 获取RPC绑定地址
func GetRPCBind() string {
	return cfg.RPC.Bind
}

// GetHttpBind 获取HTTP绑定地址
func GetHttpBind() string {
	return cfg.Http.Bind
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

// GetCertificateDir 获取根证书路径
func GetCertificateDir() string {
	return cfg.Dir.Certificate
}

// GetAliyunOss 获取阿里云OSS配置
func GetAliyunOss() *AliyunOss {
	return cfg.Aliyun.Oss
}

// GetSnca 获取SNCA配置
func GetSnca() (url, source, customerType string) {
	return cfg.Snca.Url, cfg.Snca.Source, cfg.Snca.CustomerType
}

// GetPostgresConfig 获取Postgresql配置
func GetPostgresConfig() (string, bool) {
	return cfg.Postgres.Dsn, cfg.Postgres.Debug
}

// GetConfigFile 获取配置文件路径
func GetConfigFile() string {
	return configFile
}
