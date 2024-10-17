// Copyright (C) edocseal. 2024-present.
//
// Created at 2024-04-11, by liasica

package internal

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	jsoniter "github.com/json-iterator/go"
	"github.com/spf13/cobra"

	"github.com/liasica/edocseal"
	"github.com/liasica/edocseal/internal/model"
)

func templateCommand() *cobra.Command {
	var (
		form string
		path string
	)

	cmd := &cobra.Command{
		Use:               "template <input>",
		Short:             "添加模板",
		Example:           "edocpdf template ./template.pdf",
		Args:              cobra.ExactArgs(1),
		CompletionOptions: cobra.CompletionOptions{DisableDefaultCmd: true},
		Run: func(_ *cobra.Command, args []string) {
			// 获取模板文件md5
			fileId, err := edocseal.FileMd5(args[0])
			if err != nil {
				fmt.Printf("模板MD5获取失败: %v\n", err)
				os.Exit(1)
			}

			// 重命名模板文档
			templateFile, _ := filepath.Abs(filepath.Join(path, fileId+".pdf"))
			if !edocseal.FileExists(templateFile) {
				_ = os.Rename(args[0], templateFile)
			}

			// 获取表单属性
			var b []byte
			b, err = edocseal.Exec("qpdf", form, "--json")
			if err != nil {
				fmt.Printf("模板表单解析失败: %v\n", err)
				os.Exit(1)
			}

			// 获取表单数据
			var form edocseal.Form
			err = jsoniter.Unmarshal(b, &form)
			if err != nil {
				fmt.Printf("表单数据解析失败: %v\n", err)
				os.Exit(1)
			}

			// 生成模板ID
			templateId := strings.ToUpper(strings.ReplaceAll(uuid.New().String(), "-", ""))
			template := model.Template{
				ID:     templateId,
				File:   templateFile,
				Fields: make(map[string]model.TemplateField),
			}

			// 模板数据
			var hasEnt, hasRider bool
			for _, field := range form.Acroform.Fields {
				m := form.Objects[1]["obj:"+field.Annotation.Object]
				mb, _ := jsoniter.Marshal(m)
				var data edocseal.FormFieldObject
				_ = jsoniter.Unmarshal(mb, &data)

				fmt.Printf("%20.20s %10.10s\t\t%6.2f, %6.2f, %6.2f, %6.2f\n",
					field.Fullname,
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
			_ = os.WriteFile(filepath.Join(path, fmt.Sprintf("%s.json", templateId)), json, 0755)
			fmt.Printf("完成模板创建，模板ID: %s\n", templateId)
		},
	}

	cmd.Flags().StringVarP(&form, "form", "f", "", "表单文件")
	cmd.Flags().StringVarP(&path, "path", "p", "templates", "模板路径")

	_ = cmd.MarkFlagRequired("temp")

	return cmd
}
