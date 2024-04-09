// Copyright (C) edocseal. 2024-present.
//
// Created at 2024-04-05, by liasica

package edocseal

import (
	"math/big"
	"math/rand"
	"path/filepath"

	"github.com/benoitkugler/pdf/formfill"
	"github.com/benoitkugler/pdf/model"
	"github.com/benoitkugler/pdf/reader"

	"github.com/liasica/edocseal/pb"
)

// PdfParseFields 解析PDF文件中待填充字段
func PdfParseFields(path string) (fields map[string]model.Rectangle, err error) {
	var doc model.Document
	doc, _, err = reader.ParsePDFFile(path, reader.Options{})
	if err != nil {
		return
	}

	fields = make(map[string]model.Rectangle)

	// 获取所有字段
	for _, field := range doc.Catalog.AcroForm.Flatten() {
		t := field.Field.T

		if len(field.Field.Widgets) > 0 {
			fields[t] = field.Field.Widgets[0].Rect
		}
	}

	return
}

// PdfFillForm 填充PDF表单
func PdfFillForm(filename, filledDir string, fields map[string]*pb.ContractFromField) (filled string, err error) {
	// 解析PDF文件
	var doc model.Document
	doc, _, err = reader.ParsePDFFile(filename, reader.Options{})
	if err != nil {
		return
	}

	// 获取表单字段
	var form []formfill.FDFField
	for t, value := range fields {
		field := formfill.FDFField{
			T: t,
		}
		// 根据字段类型填充值
		switch v := value.Value.(type) {
		case *pb.ContractFromField_Text:
			field.Values = formfill.Values{V: formfill.FDFText(v.Text)}
		case *pb.ContractFromField_Checkbox:
			field.Values = formfill.Values{V: formfill.FDFName(v.String())}
		}
		form = append(form, field)
	}

	form = append(form, formfill.FDFField{
		T:      "Checkbox1",
		Values: formfill.Values{V: formfill.FDFName("On")},
	})

	// 填充表单
	err = formfill.FillForm(&doc, formfill.FDFDict{Fields: form}, true)
	if err != nil {
		return
	}

	// 保存填充后的文件
	filled = filepath.Join(filledDir, FileNameWithoutExtension(filename), big.NewInt(rand.Int63()).String()+".pdf")
	err = CreateDirectory(filepath.Dir(filled))
	if err != nil {
		return
	}

	err = doc.WriteFile(filled, nil)
	return
}
