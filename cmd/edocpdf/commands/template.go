// Copyright (C) edocseal. 2024-present.
//
// Created at 2024-04-09, by liasica

package commands

import (
	"fmt"
	"os"
	"path/filepath"

	jsoniter "github.com/json-iterator/go"
	"github.com/spf13/cobra"

	"github.com/liasica/edocseal"
	"github.com/liasica/edocseal/internal/model"
)

func Template() *cobra.Command {
	var (
		temp string
		path string
	)
	cmd := &cobra.Command{
		Use:               "template <input>",
		Short:             "添加模板",
		Example:           "edocpdf template ./时光驹APP个签合同.pdf",
		Args:              cobra.ExactArgs(1),
		CompletionOptions: cobra.CompletionOptions{DisableDefaultCmd: true},
		Run: func(_ *cobra.Command, args []string) {
			// 重命名模板文档
			// 获取模板文件md5
			id, err := edocseal.FileMd5(args[0])
			if err != nil {
				fmt.Printf("模板文件MD5获取失败: %v\n", err)
				os.Exit(1)
			}
			fmt.Println("模板ID为: ", id)
			_ = os.Rename(args[0], filepath.Join(path, id+".pdf"))

			// 获取表单属性
			var b []byte
			b, err = edocseal.Exec("qpdf", temp, "--json")
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

			// 模板数据
			fields := make(map[string]model.TemplateField)
			for _, field := range form.Acroform.Fields {
				m := form.Objects[1]["obj:"+field.Annotation.Object]
				mb, _ := jsoniter.Marshal(m)
				var data edocseal.FormFieldObject
				_ = jsoniter.Unmarshal(mb, &data)

				fmt.Printf("%s (%s): [%.2f, %.2f, %.2f, %.2f]\n",
					field.Fullname,
					field.Fieldtype,
					data.Value.Rect[0],
					data.Value.Rect[1],
					data.Value.Rect[2],
					data.Value.Rect[3],
				)
				if _, ok := fields[field.Fullname]; ok {
					fmt.Printf("字段重复: %s\n", field.Fullname)
					os.Exit(1)
				}

				fields[field.Fullname] = model.TemplateField{
					Page: field.Pageposfrom1,
					Type: getFiledType(field.Fieldtype),
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
			if _, ok := fields[model.EntSignField]; !ok {
				fmt.Printf("模板没有签名字段: %s", model.EntSignField)
				os.Exit(1)
			}
			if _, ok := fields[model.PersonalSignField]; !ok {
				fmt.Printf("模板没有签名字段: %s", model.PersonalSignField)
				os.Exit(1)
			}

			// 存储模板配置
			json, _ := jsoniter.MarshalIndent(fields, "", "  ")
			_ = os.WriteFile(filepath.Join(path, fmt.Sprintf("%s.json", id)), json, 0755)
		},
	}

	cmd.Flags().StringVarP(&temp, "temp", "t", "runtime", "表单文件")
	cmd.Flags().StringVarP(&path, "path", "p", "templates", "模板路径")

	_ = cmd.MarkFlagRequired("temp")

	return cmd
}

func getFiledType(typ string) model.TemplateFieldType {
	switch typ {
	case "/Sig":
		return model.TemplateFieldTypeSignature
	case "/Tx":
		return model.TemplateFieldTypeText
	case "/Btn":
		return model.TemplateFieldTypeCheckbox
	}
	return ""
}
