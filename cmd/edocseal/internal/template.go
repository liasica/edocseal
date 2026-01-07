// Copyright (C) edocseal. 2024-present.
//
// Created at 2024-04-11, by liasica

package internal

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	jsoniter "github.com/json-iterator/go"
	"github.com/spf13/cobra"

	"auroraride.com/edocseal"
	"auroraride.com/edocseal/internal/model"
)

func generateTemplate(pdf, form, templates string) {
	// 获取模板文件md5
	templateId, err := edocseal.FileMd5(pdf)
	if err != nil {
		fmt.Printf("模板MD5获取失败: %v\n", err)
		os.Exit(1)
	}

	templateId = strings.ToUpper(templateId)

	// 重命名模板文档
	templateFile, _ := filepath.Abs(filepath.Join(templates, templateId+".pdf"))
	if !edocseal.FileExists(templateFile) {
		err = copyFileContents(pdf, templateFile)
		if err != nil {
			fmt.Printf("模板文件复制失败: %v\n", err)
			os.Exit(1)
		}
	}

	// 获取表单属性
	var b []byte
	b, err = edocseal.Exec("qpdf", form, "--json")
	if err != nil {
		fmt.Printf("模板表单解析失败: %v\n", err)
		os.Exit(1)
	}

	// 获取表单数据
	var fd edocseal.Form
	err = jsoniter.Unmarshal(b, &fd)
	if err != nil {
		fmt.Printf("表单数据解析失败: %v\n", err)
		os.Exit(1)
	}

	// 生成模板ID
	template := model.Template{
		ID:     templateId,
		File:   templateFile,
		Fields: make(map[string]model.TemplateField),
	}

	// 模板数据
	var hasEnt, hasRider bool
	for _, field := range fd.Acroform.Fields {
		m := fd.Objects[1]["obj:"+field.Annotation.Object]
		mb, _ := jsoniter.Marshal(m)
		var data edocseal.FormFieldObject
		_ = jsoniter.Unmarshal(mb, &data)

		fmt.Printf("%20.20s %20.20s\t%10.10s\t\t%6.2f, %6.2f, %6.2f, %6.2f\n",
			field.Fullname,
			field.Alternativename,
			field.Fieldtype,
			data.Value.Rect[0],
			data.Value.Rect[1],
			data.Value.Rect[2],
			data.Value.Rect[3],
		)
		// if _, ok := fields[field.Fullname]; ok {
		// 	fmt.Printf("字段重复: %s\n", field.Fullname)
		// 	os.Exit(1)
		// }

		var typ model.TemplateFieldType
		typ, err = model.NewTemplateFieldType(field.Fieldtype)
		if err != nil {
			fmt.Printf("字段类型解析失败: %v\n", err)
			os.Exit(1)
		}

		if field.Fullname == model.EntSignField {
			hasEnt = true
		}

		if field.Fullname == model.PersonalSignField {
			hasRider = true
		}

		template.Fields[field.Fullname] = model.TemplateField{
			Page: field.Pageposfrom1,
			Type: typ,
			Rectangle: &model.TemplateRectangle{
				LeftBottom: model.Coordinate{
					X: data.Value.Rect[0],
					Y: data.Value.Rect[1],
				},
				RightTop: model.Coordinate{
					X: data.Value.Rect[2],
					Y: data.Value.Rect[3],
				},
			},
		}
	}

	// 判定是否有签名字段
	if !hasEnt {
		fmt.Printf("模板没有签名字段: %s\n", model.EntSignField)
		os.Exit(1)
	}
	if !hasRider {
		fmt.Printf("模板没有签名字段: %s\n", model.PersonalSignField)
		os.Exit(1)
	}

	// 存储模板配置
	json, _ := jsoniter.MarshalIndent(template, "", "  ")
	_ = os.WriteFile(filepath.Join(templates, fmt.Sprintf("%s.json", templateId)), json, 0755)
	fmt.Printf("完成模板创建，模板ID: %s\n", templateId)
}

type Template struct {
}

func NewTemplate() *Template {
	return &Template{}
}

func (t *Template) Group() *cobra.Group {
	return &cobra.Group{
		ID:    "template",
		Title: "模板",
	}
}

func (t *Template) Command() *cobra.Command {
	cmd := &cobra.Command{
		GroupID:           t.Group().ID,
		Use:               "template",
		Short:             "模板相关命令",
		CompletionOptions: cobra.CompletionOptions{DisableDefaultCmd: true},
	}

	cmd.AddCommand(
		t.add(),
		t.scan(),
	)

	return cmd
}

func (*Template) add() *cobra.Command {
	var (
		form      string
		templates string
	)

	cmd := &cobra.Command{
		Use:               "add <input>",
		Short:             "添加模板",
		Example:           "edocseal template add ./contract.pdf --form ./form.pdf --path ./templates",
		Args:              cobra.ExactArgs(1),
		CompletionOptions: cobra.CompletionOptions{DisableDefaultCmd: true},
		Run: func(_ *cobra.Command, args []string) {
			generateTemplate(args[0], form, templates)
		},
	}

	cmd.Flags().StringVarP(&form, "form", "f", "", "表单模板")
	cmd.Flags().StringVarP(&templates, "templates", "t", "./templates", "模板路径")

	_ = cmd.MarkFlagRequired("form")

	return cmd
}

func (*Template) scan() *cobra.Command {
	var (
		templates string
		force     bool
	)

	cmd := &cobra.Command{
		Use:               "scan <input>",
		Short:             "扫描目录生成模板",
		Example:           "edocseal template scan ./runtime",
		Args:              cobra.ExactArgs(1),
		CompletionOptions: cobra.CompletionOptions{DisableDefaultCmd: true},
		Run: func(_ *cobra.Command, args []string) {
			// generateTemplate(args[0], form, templates)
			// 扫描目录
			files, err := os.ReadDir(args[0])
			if err != nil {
				fmt.Printf("读取目录失败: %v\n", err)
				os.Exit(1)
			}

			var list []string

			for _, f := range files {
				// 跳过目录、非PDF后缀文件、表单文件
				if f.IsDir() || !strings.HasSuffix(f.Name(), ".pdf") || strings.HasSuffix(f.Name(), "-form.pdf") {
					continue
				}

				// 检查是否存在表单文件
				name := strings.TrimSuffix(f.Name(), filepath.Ext(f.Name()))
				form := filepath.Join(args[0], fmt.Sprintf("%s-form.pdf", name))
				if !edocseal.FileExists(form) {
					continue
				}

				// 获取合同模板ID
				pdf := filepath.Join(args[0], f.Name())
				var templateId string
				templateId, err = edocseal.FileMd5(pdf)
				if err != nil {
					fmt.Printf("合同模板MD5获取失败: %v\n", err)
					os.Exit(1)
				}
				templateId = strings.ToUpper(templateId)

				// 检查模板文件是否存在
				templateFile := filepath.Join(templates, fmt.Sprintf("%s.pdf", templateId))
				if edocseal.FileExists(templateFile) && !force {
					fmt.Printf("模板文件已存在: %s\n", templateId)
					continue
				}

				// 生成模板
				message := fmt.Sprintf("%s: %s", name, templateId)
				fmt.Printf("生成模板: %s\n", message)
				generateTemplate(pdf, form, templates)
				list = append(list, message)

				fmt.Println("-------------------------")
			}

			fmt.Printf("完成模板扫描，共生成 %d 个模板:\n", len(list))
			for _, l := range list {
				fmt.Printf(" - %s\n", l)
			}
		},
	}

	cmd.Flags().StringVarP(&templates, "templates", "t", "./templates", "模板路径")
	cmd.Flags().BoolVarP(&force, "force", "f", false, "是否强制覆盖已存在的模板")

	return cmd
}
