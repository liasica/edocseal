// Copyright (C) edocseal. 2024-present.
//
// Created at 2024-04-11, by liasica

package biz

import (
	"path/filepath"
	"testing"
	"time"

	"auroraride.com/edocseal/internal/g"
)

func TestGetDocumentPaths(t *testing.T) {
	docId := time.Now().Format("200601") + g.GetID()
	t.Log(filepath.Join(docId[:4], docId[4:6], docId[6:], docId+"_unsigned.pdf"))
}
