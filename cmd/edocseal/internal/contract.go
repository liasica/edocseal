// Copyright (C) edocseal. 2026-present.
//
// Created at 2026-10-20, by liasica

package internal

import (
	"bufio"
	"fmt"
	"os"
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/spf13/cobra"

	"auroraride.com/edocseal/internal/biz"
	"auroraride.com/edocseal/internal/ent"
	"auroraride.com/edocseal/pb"
)

// regenerateItem 待重新生成的合同，每行一个
type regenerateItem struct {
	ContractID uint64            `json:"contractId"` // 主业务合同ID
	OldDocID   string            `json:"oldDocId"`   // 原文档编号
	TemplateID string            `json:"templateId"` // 模板编号
	Name       string            `json:"name"`       // 签署人
	Idcard     string            `json:"idcard"`     // 身份证号
	Phone      string            `json:"phone"`      // 手机号
	Province   string            `json:"province"`   // 省份
	City       string            `json:"city"`       // 城市
	Address    string            `json:"address"`    // 地址
	Image      string            `json:"image"`      // 个人签名图片（base64编码）
	Texts      map[string]string `json:"texts"`      // 文本字段
	Checks     []string          `json:"checks"`     // 勾选字段
}

// regenerateResult 重新生成结果，每行一个
type regenerateResult struct {
	ContractID uint64 `json:"contractId"`
	OldDocID   string `json:"oldDocId"`
	DocID      string `json:"docId,omitempty"`
	URL        string `json:"url,omitempty"`
	Error      string `json:"error,omitempty"`
}

func contractCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "contract",
		Short: "合同处理",
	}

	var (
		input            string
		output           string
		creditCode       string
		enterpriseIssuer string
		issuer           string
	)

	regenerate := &cobra.Command{
		Use:               "regenerate",
		Short:             "按模板重新生成合同并以指定签约企业签署，已成功的合同跳过",
		Example:           "edocseal contract regenerate --input items.jsonl --output result.jsonl --enterprise 91610133MA712R0U00 --enterprise-issuer SNCA --issuer SELF",
		CompletionOptions: cobra.CompletionOptions{DisableDefaultCmd: true},
		Run: func(_ *cobra.Command, _ []string) {
			err := regenerateContracts(input, output, creditCode, enterpriseIssuer, issuer)
			if err != nil {
				fmt.Printf("重新生成失败: %v\n", err)
				os.Exit(1)
			}
		},
	}

	regenerate.Flags().StringVar(&input, "input", "", "待重新生成的合同列表，每行一个 JSON")
	regenerate.Flags().StringVar(&output, "output", "", "结果文件，每行一个 JSON，追加写入")
	regenerate.Flags().StringVar(&creditCode, "enterprise", "", "签约企业统一社会信用代码")
	regenerate.Flags().StringVar(&enterpriseIssuer, "enterprise-issuer", "", "企业证书生成方式：SELF、SNCA")
	regenerate.Flags().StringVar(&issuer, "issuer", "", "个人证书生成方式：SELF、SNCA")
	_ = regenerate.MarkFlagRequired("input")
	_ = regenerate.MarkFlagRequired("output")
	_ = regenerate.MarkFlagRequired("enterprise")
	_ = regenerate.MarkFlagRequired("enterprise-issuer")
	_ = regenerate.MarkFlagRequired("issuer")

	cmd.AddCommand(regenerate)

	return cmd
}

func regenerateContracts(
	input string,
	output string,
	creditCode string,
	enterpriseIssuer string,
	issuer string,
) (err error) {
	// 读取已成功的合同
	done := make(map[uint64]bool)

	var items []regenerateItem

	err = readJSONLines(output, func(b []byte) error {
		var res regenerateResult
		if jsoniter.Unmarshal(b, &res) == nil && res.Error == "" && res.URL != "" {
			done[res.ContractID] = true
		}
		return nil
	})
	if err != nil && !os.IsNotExist(err) {
		return
	}

	err = readJSONLines(input, func(b []byte) error {
		var item regenerateItem
		if e := jsoniter.Unmarshal(b, &item); e != nil {
			return e
		}
		items = append(items, item)
		return nil
	})
	if err != nil {
		return
	}

	var f *os.File

	f, err = os.OpenFile(output, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return
	}

	defer func() {
		_ = f.Close()
	}()

	for i, item := range items {
		if done[item.ContractID] {
			continue
		}

		res := regenerateResult{ContractID: item.ContractID, OldDocID: item.OldDocID}

		var doc *ent.Document

		doc, err = biz.RegenerateDocument(creditCode, enterpriseIssuer, item.createRequest(), item.signRequest(issuer))
		if err != nil {
			res.Error = err.Error()
		} else {
			res.DocID, res.URL = doc.ID, doc.SignedURL
		}

		b, _ := jsoniter.Marshal(res)

		_, err = f.Write(append(b, '\n'))
		if err != nil {
			return
		}

		fmt.Printf("[%d/%d] %d %s %s%s\n", i+1, len(items), res.ContractID, res.DocID, res.URL, res.Error)
	}

	return
}

func (item *regenerateItem) createRequest() *pb.ContractServiceCreateRequest {
	values := make(map[string]*pb.ContractFormField, len(item.Texts)+len(item.Checks))
	for k, v := range item.Texts {
		values[k] = &pb.ContractFormField{Value: &pb.ContractFormField_Text{Text: v}}
	}

	for _, k := range item.Checks {
		values[k] = &pb.ContractFormField{Value: &pb.ContractFormField_Checkbox{Checkbox: true}}
	}

	return &pb.ContractServiceCreateRequest{
		TemplateId: item.TemplateID,
		Values:     values,
		Idcard:     item.Idcard,
		Expire:     time.Now().Add(time.Hour).Unix(),
	}
}

func (item *regenerateItem) signRequest(issuer string) *pb.ContractServiceSignRequest {
	return &pb.ContractServiceSignRequest{
		Image:             item.Image,
		Name:              item.Name,
		Province:          item.Province,
		City:              item.City,
		Address:           item.Address,
		Phone:             item.Phone,
		Idcard:            item.Idcard,
		CertificateIssuer: issuer,
	}
}

// readJSONLines 逐行读取文件
func readJSONLines(path string, fn func(b []byte) error) (err error) {
	var f *os.File

	f, err = os.Open(path)
	if err != nil {
		return
	}

	defer func() {
		_ = f.Close()
	}()

	scanner := bufio.NewScanner(f)
	scanner.Buffer(make([]byte, 0, 1024*1024), 16*1024*1024)

	for scanner.Scan() {
		if len(scanner.Bytes()) == 0 {
			continue
		}

		err = fn(scanner.Bytes())
		if err != nil {
			return
		}
	}

	err = scanner.Err()
	return
}
