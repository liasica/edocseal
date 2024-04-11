// Copyright (C) edocseal. 2024-present.
//
// Created at 2024-04-05, by liasica

package internal

import (
	"bufio"
	"crypto/x509"
	"crypto/x509/pkix"
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/liasica/edocseal"
	"github.com/liasica/edocseal/ca"
)

type Certificate struct{}

func certificateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "certificate",
		Short: "密钥和证书管理",
	}

	cert := new(Certificate)
	cmd.AddCommand(cert.root())
	cmd.AddCommand(cert.inert())
	return cmd
}

func (*Certificate) root() *cobra.Command {
	var (
		path     string
		override bool
	)

	cmd := &cobra.Command{
		Use:               "root",
		Short:             "生成根密钥和证书",
		CompletionOptions: cobra.CompletionOptions{DisableDefaultCmd: true},
		Run: func(cmd *cobra.Command, args []string) {
			// 根证书私钥
			privateKeyPath := filepath.Join(path, "rootPrivateKey.pem")
			// 根证书
			rootCertPath := filepath.Join(path, "rootCA.crt")

			// 判断是否存在密钥和证书
			if !override && (edocseal.FileExists(privateKeyPath) || edocseal.FileExists(rootCertPath)) {
				reader := bufio.NewReader(os.Stdin)
				fmt.Print("密钥和证书已存在，是否覆盖？(Y/n): ")

				input, _ := reader.ReadString('\n')
				input = strings.TrimSpace(input)
				if input == "" {
					input = "Y"
				}
				override = strings.ToUpper(input) == "Y"

				if !override {
					fmt.Println("密钥和证书已存在，如需重新生成请使用 -o 参数")
					os.Exit(1)
				}
			}

			// 生成密钥并保存
			priKey := ca.GenerateRsaPrivateKey()
			k, _ := x509.MarshalPKCS8PrivateKey(priKey)
			err := ca.SaveToFile(privateKeyPath, k, ca.BlocTypePrivateKey)
			if err != nil {
				fmt.Printf("密钥保存失败：%s", err)
				os.Exit(1)
			}

			// 生成根证书并保存
			var rootCertificate []byte
			rootCertificate, err = ca.GenerateRootCertificate(priKey, pkix.Name{
				Country:            []string{"中国"},                                    // 国家
				Province:           []string{"北京市"},                                   // 省份
				Locality:           []string{"东城区"},                                   // 城市
				Organization:       []string{"Tianjin NonEntia Technology Co., Ltd."}, // 证书持有者组织名称
				OrganizationalUnit: []string{"NonEntiaLtd"},                           // 证书持有者组织唯一标识
				CommonName:         "NonEntia Root CA",                                // 证书持有者通用名，需保持唯一，否则验证会失败
			})
			if err != nil {
				fmt.Printf("生成根证书失败：%s", err)
				os.Exit(1)
			}
			err = ca.SaveToFile(rootCertPath, rootCertificate, ca.BlocTypeCertificate)
			if err != nil {
				fmt.Printf("根证书保存失败：%s", err)
				os.Exit(1)
			}
		},
	}

	cmd.Flags().StringVarP(&path, "path", "p", "config", "生成的密钥和证书存放路径")
	cmd.Flags().BoolVarP(&override, "override", "o", false, "是否覆盖已有密钥和证书")

	return cmd
}

func (*Certificate) inert() *cobra.Command {
	var (
		path               string
		province           string
		locality           string
		organization       string
		organizationalUnit string
		commonName         string
		dns                string
	)

	cmd := &cobra.Command{
		Use:               "inter",
		Short:             "签发中间证书",
		CompletionOptions: cobra.CompletionOptions{DisableDefaultCmd: true},
		Run: func(cmd *cobra.Command, args []string) {
			// 根证书私钥
			rootPrivateKeyPath := filepath.Join(path, "rootPrivateKey.pem")
			// 根证书
			rootCertPath := filepath.Join(path, "rootCA.crt")

			// 判断是否存在密钥和证书
			if !edocseal.FileExists(rootPrivateKeyPath) || !edocseal.FileExists(rootCertPath) {
				fmt.Println("根证书不存在，请先执行 certificate root 命令生成根证书")
				os.Exit(1)
			}

			// 加载根证书和私钥
			rootPrikey, err := ca.LoadPrivateKeyFromFile(rootPrivateKeyPath)
			if err != nil {
				fmt.Printf("加载根证书私钥失败：%s", err)
				os.Exit(1)
			}

			var rootCa *x509.Certificate
			rootCa, err = ca.LoadCertificateFromFile(rootCertPath)
			if err != nil {
				fmt.Printf("加载根证书失败：%s", err)
				os.Exit(1)
			}

			// 签发证书并保存
			var (
				interCrt, interKey []byte
				serial             *big.Int
			)
			interCrt, interKey, serial, err = ca.CreateInterCertificate(rootPrikey, rootCa, pkix.Name{
				Country:            []string{"中国"},               // 国家
				Province:           []string{province},           // 省份
				Locality:           []string{locality},           // 城市
				Organization:       []string{organization},       // 证书持有者组织名称
				OrganizationalUnit: []string{organizationalUnit}, // 证书持有者组织唯一标识
				CommonName:         commonName,                   // 证书持有者通用名，需保持唯一，否则验证会失败
			})
			if err != nil {
				fmt.Printf("签发证书失败：%s", err)
				return
			}

			sn := serial.String()

			// 保存密钥
			err = ca.SaveToFile(filepath.Join(path, sn+".pem"), interKey, ca.BlocTypePrivateKey)
			if err != nil {
				fmt.Printf("密钥保存失败：%s", err)
				os.Exit(1)
			}

			// 保存证书
			err = ca.SaveToFile(filepath.Join(path, sn+".crt"), interCrt, ca.BlocTypeCertificate)
			if err != nil {
				fmt.Printf("证书保存失败：%s", err)
				os.Exit(1)
			}

			fmt.Printf("证书签发成功，证书序列号：%s\n", sn)
		},
	}

	cmd.Flags().StringVarP(&path, "path", "p", "config", "生成的密钥和证书存放路径")
	cmd.Flags().StringVar(&province, "province", "北京市", "省份")
	cmd.Flags().StringVar(&locality, "locality", "东城区", "城市")
	cmd.Flags().StringVar(&organization, "organization", "张三", "证书持有者组织名称")
	cmd.Flags().StringVar(&organizationalUnit, "organizationalUnit", "ZhangSan", "证书持有者组织唯一标识")
	cmd.Flags().StringVar(&commonName, "commonName", "ZhangSan CA", "证书持有者通用名")
	cmd.Flags().StringVar(&dns, "dns", "liasica.com", "证书域名")

	return cmd
}
