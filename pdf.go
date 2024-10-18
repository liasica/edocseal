// Copyright (C) edocseal. 2024-present.
//
// Created at 2024-04-05, by liasica

package edocseal

import (
	"bytes"
	"io"
	"os"

	"github.com/signintech/gopdf"
)

const DefaultMediaBox = "/MediaBox"

type PdfCreator struct {
	box   string
	fonts map[string]io.Reader
	font  *PdfFont
}

type PdfFont struct {
	Name  string
	Size  any
	Style string
}

type PdfCreateOption interface {
	apply(*PdfCreator)
}

type pdfCreateOptionFunc func(*PdfCreator)

func (f pdfCreateOptionFunc) apply(creator *PdfCreator) {
	f(creator)
}

func WithMediaBox(mediaBox string) PdfCreateOption {
	return pdfCreateOptionFunc(func(creator *PdfCreator) {
		creator.box = mediaBox
	})
}

func WithFont(name string, font io.Reader) PdfCreateOption {
	return pdfCreateOptionFunc(func(creator *PdfCreator) {
		if creator.fonts == nil {
			creator.fonts = make(map[string]io.Reader)
		}
		creator.fonts[name] = font
	})
}

func WithDefaultFont(name string, style string, size any) PdfCreateOption {
	return pdfCreateOptionFunc(func(creator *PdfCreator) {
		creator.font = &PdfFont{
			Name:  name,
			Size:  size,
			Style: style,
		}
	})
}

func NewPdfCreator(opts ...PdfCreateOption) *PdfCreator {
	c := &PdfCreator{
		box: DefaultMediaBox,
	}
	for _, opt := range opts {
		opt.apply(c)
	}
	return c
}

// CreatePDF 创建并处理PDF
func (creator *PdfCreator) CreatePDF(out string, source io.ReadSeeker, process func(*gopdf.GoPdf) error) (b []byte, err error) {
	// 创建PDF
	pdf := &gopdf.GoPdf{}
	pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4})

	// 导入页面
	err = pdf.ImportPagesFromSource(source, creator.box)
	if err != nil {
		return
	}

	// 加载字体
	for name, font := range creator.fonts {
		err = pdf.AddTTFFontByReader(name, font)
		if err != nil {
			return
		}
	}

	// 设置默认字体
	if creator.font != nil {
		err = pdf.SetFont(creator.font.Name, creator.font.Style, creator.font.Size)
		if err != nil {
			return
		}
	}

	// 处理PDF
	err = process(pdf)
	if err != nil {
		return
	}

	// 文档写入缓冲区
	buf := new(bytes.Buffer)
	_, err = pdf.WriteTo(buf)
	if err != nil {
		return
	}

	// 保存文档
	b = buf.Bytes()
	err = os.WriteFile(out, b, os.ModePerm)
	return
}
