package internal

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"auroraride.com/edocseal/internal/biz"
)

func enterpriseCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "enterprise",
		Short: "签约企业查询，新增、修改、切换与作废证书在管理后台操作",
	}

	cmd.AddCommand(&cobra.Command{
		Use:               "list",
		Short:             "列出签约企业与企业证书",
		CompletionOptions: cobra.CompletionOptions{DisableDefaultCmd: true},
		Run: func(_ *cobra.Command, _ []string) {
			for _, item := range biz.ListEnterprises() {
				fmt.Printf("%s %s 签约企业=%t 签章=%t\n", item.CreditCode, item.Name, item.IsDefault, len(item.Seal) > 0)
				for _, cert := range item.Certificates {
					fmt.Printf("  %s %s 颁发者 %s 到期 %s\n", cert.Issuer, cert.Serial, cert.IssuerName, time.Unix(cert.NotAfter, 0).Format(time.DateTime))
				}
			}
		},
	})

	return cmd
}
