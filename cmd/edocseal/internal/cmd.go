// Copyright (C) edocseal. 2024-present.
//
// Created at 2024-04-05, by liasica

package internal

import (
	"github.com/spf13/cobra"

	"auroraride.com/edocseal/internal"
	"auroraride.com/edocseal/internal/g"
)

func RunCommand() {
	var configFile string

	cmd := &cobra.Command{
		Use:               "edocseal",
		Short:             "电子签名控制台",
		CompletionOptions: cobra.CompletionOptions{DisableDefaultCmd: true},
		PersistentPreRun: func(_ *cobra.Command, _ []string) {
			// 加载配置文件
			g.LoadConfig(configFile)

			// 初始化
			internal.Boot()
		},
	}

	cmd.PersistentFlags().StringVarP(&configFile, "config", "c", "config/config.yaml", "配置文件")

	cmd.AddGroup(&cobra.Group{
		ID:    "template",
		Title: "合同模板相关命令",
	})

	tmpl := NewTemplate()
	cmd.AddGroup(tmpl.Group())

	cmd.AddCommand(
		certificateCommand(),
		serverCommand(),
		tmpl.Command(),
		shorturlCommand(),
		enterpriseCommand(),
	)

	_ = cmd.Execute()
}
