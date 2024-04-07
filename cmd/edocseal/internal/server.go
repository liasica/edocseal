// Copyright (C) edocseal. 2024-present.
//
// Created at 2024-04-07, by liasica

package internal

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/liasica/edocseal/internal/g"
)

func serverCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:               "server",
		Short:             "启动服务端",
		CompletionOptions: cobra.CompletionOptions{DisableDefaultCmd: true},
		Run: func(_ *cobra.Command, _ []string) {
			// 判定企业证书是否设置
			_, entCrt := g.NewCertificate()
			if !entCrt.IsValid() {
				fmt.Println("企业证书无效")
				os.Exit(1)
			}

			// TODO: 其他服务端启动逻辑
			select {}
		},
	}

	return cmd
}
