// Copyright (C) edocseal. 2024-present.
//
// Created at 2024-04-08, by liasica

package edocseal

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	"os"
	"testing"

	"github.com/benoitkugler/pdf/model"
	"github.com/benoitkugler/pdf/reader"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu"
	"github.com/signintech/gopdf"
	"github.com/signintech/pdft"
	xpdf "github.com/signintech/pdft/minigopdf"
	"github.com/stretchr/testify/require"

	"github.com/liasica/edocseal/pb"
)

func TestPdfParseFields(t *testing.T) {
	fields, err := PdfParseFields("templates/input.pdf")
	require.NoError(t, err)
	t.Log(fields)
}

func TestGetFormField(t *testing.T) {
	p := "./runtime/x.pdf"
	doc, _, err := reader.ParsePDFFile(p, reader.Options{})
	require.NoError(t, err)

	// form := doc.Catalog.AcroForm
	// fmt.Println(form)

	// 获取所有字段
	var f1, f2 model.FormFieldInherited
	for _, field := range doc.Catalog.AcroForm.Flatten() {
		f := field.Field
		w := field.Field.Widgets[0]
		if f.T == "fill_1" {
			f1 = field
			rect := w.Rect
			// fmt.Println(f, w, rect)
			fmt.Printf("%#v", rect)
		}
		if f.T == "fill_1_2" {
			f2 = field
			rect := w.Rect
			// fmt.Println(f, w, rect)
			fmt.Printf("%#v", rect)
		}
	}

	fmt.Println(f1, f2)

	// form.Fields = nil
	// err = doc.WriteFile("runtime/x.output.pdf", nil)
	// require.NoError(t, err)
	// var ctx *model.Context
	// ctx, err = pdfcpu.ReadFile(p, nil)
	// require.NoError(t, err)
	//
	// fmt.Println(ctx)
}

func TestPdfRead(t *testing.T) {
	p := "./runtime/x.pdf"
	b, err := os.ReadFile(p)
	require.NoError(t, err)

	fmt.Println(b)
}

func TestPdfFillForm(t *testing.T) {
	filled, err := PdfFillForm("runtime/input-s.pdf", "runtime", map[string]*pb.ContractFromField{
		"toggle_1": {
			Value: &pb.ContractFromField_Checkbox{
				Checkbox: true,
			},
		},
		"fill_1": {
			Value: &pb.ContractFromField_Text{
				Text: "Hello, World!",
			},
		},
	})
	require.NoError(t, err)
	t.Log(filled)
}

func TestPdft(t *testing.T) {
	pt := new(pdft.PDFt)
	err := pt.Open("./runtime/template.pdf")
	require.NoError(t, err)

	err = pt.AddFont("Song", "./runtime/HuawenFangSong.ttf")
	require.NoError(t, err)

	err = pt.SetFont("Song", "", 10)
	require.NoError(t, err)

	// lower-left x, lower-left y, upper-right x, and upper-right y
	// model.Rectangle{Llx:144.24, Lly:385.68, Urx:307.08, Ury:404.64}
	// 144.24, 385.68, 307.08, 404.64
	// A4大小 W: 595, H: 842 (DP: 72)
	// PDF编码规则是从左往右从下往上，因此此处需要使用高度 - Y
	err = pt.Insert("我是一大堆车辆型号啊", 1, 144.24, 842-404.64, 163, 19, xpdf.Left|xpdf.Middle, nil)
	require.NoError(t, err)

	var check []byte
	check, err = os.ReadFile("./internal/g/assets/check.png")
	require.NoError(t, err)

	err = pt.InsertImg(check, 1, 123.6, 842-414.904, 10, 10)
	require.NoError(t, err)

	err = pt.Save("./runtime/example_pdft.pdf")
	require.NoError(t, err)
}

func TestGoPdf(t *testing.T) {
	b, err := os.ReadFile("./internal/g/assets/check.png")
	require.NoError(t, err)

	var check image.Image
	check, err = png.Decode(bytes.NewReader(b))
	require.NoError(t, err)

	pdf := gopdf.GoPdf{}
	pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4})

	err = pdf.ImportPagesFromFile("./runtime/template.pdf", "/MediaBox")
	require.NoError(t, err)

	err = pdf.AddTTFFont("Song", "./runtime/HuawenFangSong.ttf")
	require.NoError(t, err)

	err = pdf.SetFont("Song", "", 10)
	require.NoError(t, err)

	err = pdf.SetPage(1)
	require.NoError(t, err)

	pdf.SetXY(144.24, 842-404.64)
	err = pdf.CellWithOption(&gopdf.Rect{W: 163, H: 19}, "我是一大堆车辆型号啊", gopdf.CellOption{Align: gopdf.Left | gopdf.Middle})
	require.NoError(t, err)

	err = pdf.ImageFrom(check, 123.6, 842-414.904, &gopdf.Rect{W: 10, H: 10})
	require.NoError(t, err)

	err = pdf.WritePdf("./runtime/example_gopdf.pdf")
	require.NoError(t, err)
}

func TestPdfCPU(t *testing.T) {
	ctx, err := pdfcpu.ReadFile("./runtime/template.pdf", nil)
	require.NoError(t, err)

	ref := ctx.Info
	t.Log(ref)
}
