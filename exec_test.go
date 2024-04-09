// Copyright (C) edocseal. 2024-present.
//
// Created at 2024-04-09, by liasica

package edocseal

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestExec(t *testing.T) {
	out, err := Exec("ls", "-lah")
	require.NoError(t, err)

	t.Log(string(out))
}
