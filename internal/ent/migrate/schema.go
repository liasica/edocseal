// Code generated by ent, DO NOT EDIT.

package migrate

import (
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/dialect/sql/schema"
	"entgo.io/ent/schema/field"
)

var (
	// CertificationColumns holds the columns for the "certification" table.
	CertificationColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt, Increment: true},
		{Name: "id_card_number", Type: field.TypeString, Unique: true, Size: 18, Comment: "身份证号码"},
		{Name: "private_path", Type: field.TypeString, Size: 255, Comment: "私钥路径"},
		{Name: "cert_path", Type: field.TypeString, Size: 255, Comment: "证书路径"},
		{Name: "expires_at", Type: field.TypeTime, Comment: "证书过期时间"},
	}
	// CertificationTable holds the schema information for the "certification" table.
	CertificationTable = &schema.Table{
		Name:       "certification",
		Columns:    CertificationColumns,
		PrimaryKey: []*schema.Column{CertificationColumns[0]},
	}
	// DocumentColumns holds the columns for the "document" table.
	DocumentColumns = []*schema.Column{
		{Name: "id", Type: field.TypeString, Unique: true, Comment: "文档ID"},
		{Name: "hash", Type: field.TypeString, Comment: "参数哈希"},
		{Name: "status", Type: field.TypeEnum, Comment: "文档状态，文档超时后需删除", Enums: []string{"unsigned", "signed"}, Default: "unsigned"},
		{Name: "template_id", Type: field.TypeString, Comment: "模板ID"},
		{Name: "id_card_number", Type: field.TypeString, Size: 18, Comment: "身份证号码"},
		{Name: "expires_at", Type: field.TypeTime, Comment: "过期时间"},
		{Name: "signed_url", Type: field.TypeString, Nullable: true, Comment: "已签约短链接"},
		{Name: "unsigned_url", Type: field.TypeString, Nullable: true, Comment: "已签约短链接"},
		{Name: "paths", Type: field.TypeJSON, Comment: "文档各项路径"},
		{Name: "create_at", Type: field.TypeTime, Comment: "创建时间"},
	}
	// DocumentTable holds the schema information for the "document" table.
	DocumentTable = &schema.Table{
		Name:       "document",
		Columns:    DocumentColumns,
		PrimaryKey: []*schema.Column{DocumentColumns[0]},
		Indexes: []*schema.Index{
			{
				Name:    "document_hash",
				Unique:  false,
				Columns: []*schema.Column{DocumentColumns[1]},
			},
			{
				Name:    "document_template_id",
				Unique:  false,
				Columns: []*schema.Column{DocumentColumns[3]},
			},
			{
				Name:    "document_id_card_number",
				Unique:  false,
				Columns: []*schema.Column{DocumentColumns[4]},
			},
			{
				Name:    "document_expires_at",
				Unique:  false,
				Columns: []*schema.Column{DocumentColumns[5]},
			},
		},
	}
	// Tables holds all the tables in the schema.
	Tables = []*schema.Table{
		CertificationTable,
		DocumentTable,
	}
)

func init() {
	CertificationTable.Annotation = &entsql.Annotation{
		Table: "certification",
	}
	DocumentTable.Annotation = &entsql.Annotation{
		Table: "document",
	}
}
