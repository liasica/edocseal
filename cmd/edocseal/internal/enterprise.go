package internal

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"

	"auroraride.com/edocseal/internal/biz"
	"auroraride.com/edocseal/internal/ent"
)

type Enterprise struct{}

func enterpriseCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "enterprise",
		Short: "签约企业管理",
	}

	ep := new(Enterprise)
	cmd.AddCommand(ep.save(), ep.list(), ep.revoke())
	return cmd
}

// 新增或更新签约企业
func (*Enterprise) save() *cobra.Command {
	e := new(ent.Enterprise)

	cmd := &cobra.Command{
		Use:               "save",
		Short:             "新增或更新签约企业，签章文件需按统一社会信用代码命名放入签章目录",
		CompletionOptions: cobra.CompletionOptions{DisableDefaultCmd: true},
		Run: func(_ *cobra.Command, _ []string) {
			if e.CreditCode == "" || e.Name == "" {
				fmt.Println("--code 与 --name 不能为空")
				os.Exit(1)
			}

			err := biz.SaveEnterprise(e)
			if err != nil {
				fmt.Printf("保存企业失败：%s\n", err)
				os.Exit(1)
			}

			fmt.Printf("企业已保存：%s %s\n", e.CreditCode, e.Name)
		},
	}

	cmd.Flags().StringVar(&e.CreditCode, "code", "", "统一社会信用代码")
	cmd.Flags().StringVar(&e.Name, "name", "", "企业名称")
	cmd.Flags().StringVar(&e.Province, "province", "陕西省", "省份")
	cmd.Flags().StringVar(&e.City, "city", "西安市", "城市")
	cmd.Flags().StringVar(&e.PersonName, "person", "", "代办人姓名，向陕西CA 申请证书时使用")
	cmd.Flags().StringVar(&e.Phone, "phone", "", "代办人手机号，向陕西CA 申请证书时使用")
	cmd.Flags().StringVar(&e.Idcard, "idcard", "", "代办人身份证号，向陕西CA 申请证书时使用")
	cmd.Flags().BoolVar(&e.IsDefault, "default", false, "设为默认企业")

	return cmd
}

// 列出签约企业与企业证书
func (*Enterprise) list() *cobra.Command {
	return &cobra.Command{
		Use:               "list",
		Short:             "列出签约企业与企业证书",
		CompletionOptions: cobra.CompletionOptions{DisableDefaultCmd: true},
		Run: func(_ *cobra.Command, _ []string) {
			items, certs, err := biz.ListEnterprises()
			if err != nil {
				fmt.Printf("查询企业失败：%s\n", err)
				os.Exit(1)
			}

			for _, item := range items {
				fmt.Printf("%s %s 默认=%t\n", item.CreditCode, item.Name, item.IsDefault)
				for _, cert := range certs {
					if cert.CreditCode == item.CreditCode {
						fmt.Printf("  %s 到期 %s %s\n", cert.Issuer, cert.ExpiresAt.Local().Format(time.DateTime), cert.CertPath)
					}
				}
			}
		},
	}
}

// 作废企业证书
func (*Enterprise) revoke() *cobra.Command {
	var creditCode string

	cmd := &cobra.Command{
		Use:               "revoke",
		Short:             "作废企业的全部证书，下次签约时重新签发",
		CompletionOptions: cobra.CompletionOptions{DisableDefaultCmd: true},
		Run: func(_ *cobra.Command, _ []string) {
			if creditCode == "" {
				fmt.Println("--code 不能为空")
				os.Exit(1)
			}

			n, err := biz.RevokeEnterpriseCertificates(creditCode)
			if err != nil {
				fmt.Printf("作废企业证书失败：%s\n", err)
				os.Exit(1)
			}

			fmt.Printf("已作废 %d 条企业证书记录\n", n)
		},
	}

	cmd.Flags().StringVar(&creditCode, "code", "", "统一社会信用代码")

	return cmd
}
