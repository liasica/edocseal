// Copyright (C) edocseal. 2024-present.
//
// Created at 2024-04-15, by liasica

package internal

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/liasica/edocseal/internal/biz"
)

func shorturlCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:               "shorturl <url>",
		Short:             "生成短链接",
		Example:           "edocseal shorturl https://www.baidu.com",
		Args:              cobra.ExactArgs(1),
		CompletionOptions: cobra.CompletionOptions{DisableDefaultCmd: true},
		Run: func(cmd *cobra.Command, args []string) {
			// 生成短链接
			shortUrl, err := biz.CreateShortUrl(args[0])
			if err != nil {
				fmt.Printf("生成短链接失败: %v\n", err)
				os.Exit(1)
			}
			cmd.Println(shortUrl)
		},
	}

	return cmd
}
