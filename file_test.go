// Copyright (C) edocseal. 2024-present.
//
// Created at 2024-04-08, by liasica

package edocseal

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDownloadFile(t *testing.T) {
	path, sn, err := DownloadFile("https://cdn.auroraride.com/%E5%88%98%E6%B8%8A-140603199808201635/contracts/%E6%97%B6%E5%85%89%E9%A9%B9%E9%AA%91%E6%89%8B%E7%A7%9F%E8%B5%81%E5%90%88%E5%90%8C-20240405121337252.pdf", "templates")
	require.NoError(t, err)
	t.Log(sn)
	t.Log(path)
}
