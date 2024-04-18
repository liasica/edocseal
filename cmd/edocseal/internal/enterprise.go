package internal

import (
	"bufio"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"errors"
	"fmt"
	"os"
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
		keyPath  string // 根证书私钥路径
		crtPath  string // 根证书路径
		override bool   // 是否重写
		orgName  string // 组织名称
		perName  string // 代理人名称
		phone    string // 代理人手机号码
		idCard   string // 代理人身份证号
		province string // 省份
		city     string // 城市
		orgCode  string // 企业信用代码
	)

	cmd := &cobra.Command{
		Use:               "ca",
		Short:             "生成CA企业级证书",
		CompletionOptions: cobra.CompletionOptions{DisableDefaultCmd: true},
		Run: func(cmd *cobra.Command, args []string) {
			var err error

			// 判断参数必填项
			err = checkParam(keyPath, crtPath, orgName, perName, phone, idCard, province, city, orgCode)
			if err != nil {
				fmt.Println("参数检验失败：", err)
				os.Exit(1)
			}

			// 判断是否存在密钥和证书
			checExist(keyPath, crtPath, override)

			// 生成秘钥并保存
			var priKey *rsa.PrivateKey
			priKey, err = generateEnterprisePrivatekey(keyPath)
			if err != nil {
				fmt.Println("生成企业证书秘钥失败：", err)
				os.Exit(1)
			}

			// 生成cert证书并保存
			err = generateEnterpriseCert(crtPath, orgName, perName, phone, idCard, province, city, orgCode, priKey)
			if err != nil {
				fmt.Println("生成企业证书cert失败：", err)
				os.Exit(1)
			}

		},
	}

	cmd.Flags().StringVarP(&keyPath, "keypath", "", "", "生成的企业证书私钥存放路径")
	cmd.Flags().StringVarP(&crtPath, "crtpath", "", "", "生成的企业证书Crt存放路径")
	cmd.Flags().BoolVarP(&override, "override", "", false, "是否覆盖已有企业密钥和企业证书")
	cmd.Flags().StringVarP(&orgName, "oname", "", "", "企业证书申请参数-组织名称")
	cmd.Flags().StringVarP(&perName, "pname", "", "", "企业证书申请参数-代理人名称")
	cmd.Flags().StringVarP(&phone, "phone", "", "", "企业证书申请参数-代理人手机号码")
	cmd.Flags().StringVarP(&idCard, "idcard", "", "", "企业证书申请参数-代理人身份证号")
	cmd.Flags().StringVarP(&province, "province", "", "", "企业证书申请参数-省份")
	cmd.Flags().StringVarP(&city, "city", "", "", "企业证书申请参数-城市")
	cmd.Flags().StringVarP(&orgCode, "orgcode", "", "", "企业证书申请参数-企业信用代码")

	return cmd
}

func checkParam(kPath, cPath, oName, pName, phone, idcard, provice, city, oCode string) (err error) {
	if kPath == "" {
		return errors.New("企业证书申请参数私钥存放路径 --keypath 不能为空")
	}
	if cPath == "" {
		return errors.New("企业证书申请参数Crt证书存放路径 --crtpath 不能为空")
	}
	if oName == "" {
		return errors.New("企业证书申请参数组织名称 --oname 不能为空")
	}
	if pName == "" {
		return errors.New("企业证书申请参数代理人名称 --pname 不能为空")
	}
	if phone == "" {
		return errors.New("企业证书申请参数代理人手机号码 --phone 不能为空")
	}
	if idcard == "" {
		return errors.New("企业证书申请参数代理人身份证号 --idcard 不能为空")
	}
	if provice == "" {
		return errors.New("企业证书申请参数省份 --province 不能为空")
	}
	if city == "" {
		return errors.New("企业证书申请参数城市 --city 不能为空")
	}
	if oCode == "" {
		return errors.New("企业证书申请参数企业信用代码 --orgcode 不能为空")
	}

	return
}

func checExist(kPath, cPath string, override bool) {
	if !override && (edocseal.FileExists(kPath) || edocseal.FileExists(cPath)) {
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

func generateEnterpriseCert(cPath, oName, pName, phone, idcard, provice, city, oCode string, priKey *rsa.PrivateKey) (err error) {
	// 生成CSR请求
	var csr []byte
	csr, err = ca.GenerateRequest(priKey, pkix.Name{
		Country:    []string{"CN"},
		CommonName: oName,
	})
	if err != nil {
		return
	}

	// 请求第三方申请企业证书
	var cert []byte
	cert, err = snca.NewSnca(g.GetSnca()).RequestCACert(
		snca.CertTypeEnterprise,
		oName,
		pName,
		phone,
		idcard,
		provice,
		city,
		csr,
		oCode,
	)
	if err != nil {
		return
	}
	return ca.SaveToFile(cPath, cert, ca.BlocTypeCertificate)
}
