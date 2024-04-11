// Copyright (C) edocseal. 2024-present.
//
// Created at 2024-04-09, by liasica

package edocseal

import (
	"bytes"
	"os/exec"
	"strings"
)

// Exec 执行shell命令
func Exec(arg ...string) ([]byte, error) {
	cmd := exec.Command("bash", "-c", strings.Join(arg, " "))
	stdout, err := cmd.Output()

	if err != nil {
		return nil, err
	}

	stdout = bytes.TrimRight(stdout, "\n")

	return stdout, nil
}
