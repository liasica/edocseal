// Copyright (C) edocseal. 2024-present.
//
// Created at 2024-04-05, by liasica

package edocseal

import (
	"log"

	"github.com/benoitkugler/pdf/formfill"
	"github.com/benoitkugler/pdf/reader"
)

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
