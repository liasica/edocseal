// Copyright (C) edocseal. 2024-present.
//
// Created at 2024-04-11, by liasica

package biz

import (
	"context"
	"path/filepath"

	"auroraride.com/edocseal/internal/ent"
	"auroraride.com/edocseal/internal/ent/document"
	"auroraride.com/edocseal/internal/g"
	"auroraride.com/edocseal/internal/model"
)

func QueryDocument(docId string) (*ent.Document, error) {
	// 查询文档
	return ent.NewDatabase().Document.Query().
		Where(document.ID(docId)).
		First(context.Background())
}

func NewPaths(docId string) *model.Paths {
	prefix := docId[:4] + "/" + docId[4:6] + "/" + docId[6:8] + "/" + docId[8:] + "/"
	return &model.Paths{
		UnSigned:    filepath.Join(g.GetDocumentDir(), prefix+"unsigned.pdf"),
		Signed:      filepath.Join(g.GetDocumentDir(), prefix+"signed.pdf"),
		Image:       filepath.Join(g.GetDocumentDir(), prefix+"/image.png"),
		OssUnSigned: "__contracts/" + prefix + "unsigned.pdf",
		OssSigned:   "__contracts/" + prefix + "signed.pdf",
	}
}
