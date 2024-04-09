// Copyright (C) edocseal. 2024-present.
//
// Created at 2024-04-09, by liasica

package main

import (
	"github.com/spf13/cobra"

	"github.com/liasica/edocseal/cmd/edocpdf/commands"
)

func main() {
	cmd := &cobra.Command{
		Use:               "edocpdf",
		Short:             "PDF工具",
		CompletionOptions: cobra.CompletionOptions{DisableDefaultCmd: true},
		Run:               func(_ *cobra.Command, _ []string) {},
	}

	cmd.AddCommand(commands.Template())

	_ = cmd.Execute()
}
