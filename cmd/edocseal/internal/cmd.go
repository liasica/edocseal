// Copyright (C) edocseal. 2024-present.
//
// Created at 2024-04-05, by liasica

package internal

import "github.com/spf13/cobra"

func RunCommand() {
	cmd := &cobra.Command{
		Use:               "edocseal",
		Short:             "电子签名控制台",
		CompletionOptions: cobra.CompletionOptions{DisableDefaultCmd: true},
		PersistentPreRun: func(_ *cobra.Command, _ []string) {
		},
	}

	cmd.AddCommand(
		certificateCommand(),
	)

	_ = cmd.Execute()
}
