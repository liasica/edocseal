// Copyright (C) edocseal. 2024-present.
//
// Created at 2024-04-07, by liasica

package internal

import (
	"github.com/spf13/cobra"
)

func serverCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:               "server",
		Short:             "启动服务端",
		CompletionOptions: cobra.CompletionOptions{DisableDefaultCmd: true},
		Run: func(_ *cobra.Command, _ []string) {
			// TODO: 其他服务端启动逻辑
			select {}
		},
	}

	return cmd
}
