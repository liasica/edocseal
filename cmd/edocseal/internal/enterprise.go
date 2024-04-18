package internal

import (
	"bufio"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/liasica/edocseal"
	"github.com/liasica/edocseal/ca"
	"github.com/liasica/edocseal/internal/g"
	"github.com/liasica/edocseal/third/snca"
)

type Enterprise struct{}

func enterpriseCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "enterprise",
		Short: "CA企业证书管理",
	}

	ep := new(Enterprise)
	cmd.AddCommand(ep.ca())
	return cmd
}

func (*Enterprise) ca() *cobra.Command {
	var (
		path     string
		override bool
	)

	cmd := &cobra.Command{
		Use:               "ca",
		Short:             "生成CA企业级证书",
		CompletionOptions: cobra.CompletionOptions{DisableDefaultCmd: true},
		Run: func(cmd *cobra.Command, args []string) {
			var err error
			// 企业根证书私钥路径
			entPKPath := filepath.Join(path, "enterprisePrivateKey.pem")
			// 企业根证书路径
			entCertPath := filepath.Join(path, "enterpriseCA.crt")

			// 判断是否存在密钥和证书
			if !override && (edocseal.FileExists(entPKPath) || edocseal.FileExists(entCertPath)) {
				reader := bufio.NewReader(os.Stdin)
				fmt.Print("企业密钥和证书已存在，是否覆盖？(Y/n): ")

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

			// 生成秘钥并保存
			var priKey *rsa.PrivateKey
			priKey, err = generateEnterprisePrivatekey(entPKPath)
			if err != nil {
				fmt.Println("生成企业证书秘钥失败：", err)
				os.Exit(1)
			}

			// 生成cert证书并保存
			err = generateEnterpriseCert(entCertPath, priKey)
			if err != nil {
				fmt.Println("生成企业证书cert失败：", err)
				os.Exit(1)
			}

		},
	}

	cmd.Flags().StringVarP(&path, "path", "p", "config", "生成的企业密钥和企业证书存放路径")
	cmd.Flags().BoolVarP(&override, "override", "o", false, "是否覆盖已有企业密钥和企业证书")

	return cmd
}

func generateEnterprisePrivatekey(p string) (priKey *rsa.PrivateKey, err error) {
	priKey = ca.GenerateRsaPrivateKey()
	var key []byte
	key, err = x509.MarshalPKCS8PrivateKey(priKey)
	if err != nil {
		return
	}
	return priKey, ca.SaveToFile(p, key, ca.BlocTypePrivateKey)
}

func generateEnterpriseCert(p string, priKey *rsa.PrivateKey) (err error) {
	name := "西安时光驹新能源科技有限公司"
	// 生成CSR请求
	var csr []byte
	csr, err = ca.GenerateRequest(priKey, pkix.Name{
		Country:    []string{"CN"},
		CommonName: name,
	})
	if err != nil {
		return
	}

	// 请求第三方申请企业证书
	var cert []byte
	cert, err = snca.NewSnca(g.GetSnca()).RequestCACert(
		snca.CertTypeEnterprise,
		name,
		"李四",
		"18555555555",
		"110101199003070000",
		"陕西省",
		"西安市",
		csr,
		"91610133MA712R0U00",
	)
	if err != nil {
		return
	}
	return ca.SaveToFile(p, cert, ca.BlocTypeCertificate)
}
