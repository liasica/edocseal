// Copyright (C) edocseal. 2024-present.
//
// Created at 2024-04-14, by liasica

package g

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetFont(t *testing.T) {
	f, err := GetFont("Song")
	require.NoError(t, err)
	t.Log(f.Stat())
}
