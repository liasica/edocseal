// Copyright (C) edocseal. 2024-present.
//
// Created at 2024-04-05, by liasica

package internal

import (
	"bufio"
	"crypto/x509"
	"crypto/x509/pkix"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/liasica/edocseal/ca"
)

func certificateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "certificate",
		Short: "密钥和证书管理",
	}

	cmd.AddCommand(certificateGenerateCommand())
	return cmd
}

func certificateGenerateCommand() *cobra.Command {
	var (
		path     string
		override bool
	)

	cmd := &cobra.Command{
		Use:   "generate",
		Short: "生成密钥和证书",
		Run: func(cmd *cobra.Command, args []string) {
			var (
				privateKeyPath = path + "/rootPrivateKey.pem" // 根证书私钥
				rootCertPath   = path + "/rootCA.crt"         // 根证书
			)

			// 判断是否存在密钥和证书
			if !override && (fileExists(privateKeyPath) || fileExists(rootCertPath)) {
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
					return
				}
			}

			// 生成密钥并保存
			priKey := ca.GenerateRsaPrivateKey()
			err := os.WriteFile(privateKeyPath, ca.PEMEncoding(x509.MarshalPKCS1PrivateKey(priKey), ca.BlocTypePrivateKey), 0644)
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
			err = os.WriteFile(rootCertPath, ca.PEMEncoding(rootCertificate, ca.BlocTypeCertificate), 0644)
			if err != nil {
				fmt.Printf("根证书保存失败：%s", err)
				os.Exit(1)
			}
		},
	}

	cmd.PersistentFlags().StringVarP(&path, "path", "p", "config", "生成的密钥和证书存放路径")
	cmd.PersistentFlags().BoolVarP(&override, "override", "o", false, "是否覆盖已有密钥和证书")

	return cmd
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
