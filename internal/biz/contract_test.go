// Copyright (C) edocseal. 2024-present.
//
// Created at 2024-04-10, by liasica

package biz

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/liasica/edocseal"
	"github.com/liasica/edocseal/internal/g"
	"github.com/liasica/edocseal/pb"
)

func TestCreateDocument(t *testing.T) {
	g.LoadConfig("config/config.yaml")
	err := edocseal.CreateDirectory(g.GetDocumentDir())
	require.NoError(t, err)

	var (
		b     []byte
		docId string
	)
	b, docId, err = CreateDocument("81f52261d44e5dfc756a86ebae720285", map[string]*pb.ContractFromField{
		"check1":       {Value: &pb.ContractFromField_Checkbox{Checkbox: true}},
		"name":         {Value: &pb.ContractFromField_Text{Text: "张三"}},
		"idcardNumber": {Value: &pb.ContractFromField_Text{Text: "110101199003070000"}},
	})
	require.NoError(t, err)
	t.Logf("文档已创建, 文档ID: %s, 大小: %dbytes", docId, len(b))
}
