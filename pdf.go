// Copyright (C) edocseal. 2024-present.
//
// Created at 2024-04-05, by liasica

package edocseal

import (
	"log"

	"github.com/benoitkugler/pdf/formfill"
	"github.com/benoitkugler/pdf/model"
	"github.com/benoitkugler/pdf/reader"
)

// ParsePdfFields 解析PDF文件中待填充字段
func ParsePdfFields(path string) (fields map[string]model.Rectangle, err error) {
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

func Sign() {
	doc, _, err := reader.ParsePDFFile("input.pdf", reader.Options{})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("got %d fields", len(doc.Catalog.AcroForm.Flatten()))

	err = formfill.FillForm(&doc, formfill.FDFDict{Fields: []formfill.FDFField{
		{
			T:      "fill_1",
			Values: formfill.Values{V: formfill.FDFText("My text sample 1")},
		},
		{
			T:      "fill_2",
			Values: formfill.Values{V: formfill.FDFText("Hello World!")},
		},
		{
			T:      "toggle_1",
			Values: formfill.Values{V: formfill.FDFName("On")},
		},
	}}, true)
	if err != nil {
		log.Fatal(err)
	}
	err = doc.WriteFile("output.pdf", nil)
	if err != nil {
		log.Fatal(err)
	}
}
