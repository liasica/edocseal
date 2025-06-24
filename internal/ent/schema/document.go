package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"

	"auroraride.com/edocseal/internal/model"
)

// Document holds the schema definition for the Document entity.
type Document struct {
	ent.Schema
}

// Annotations of the Document.
func (Document) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "document"},
		entsql.WithComments(true),
	}
}

// Fields of the Document.
func (Document) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").Unique().Comment("文档ID"),
		field.String("hash").Comment("参数哈希"),
		field.Enum("status").Default("unsigned").Values("unsigned", "signed").Comment("文档状态，文档超时后需删除"),
		field.String("template_id").Comment("模板ID"),
		field.String("id_card_number").MaxLen(18).Comment("身份证号码"),
		field.Time("expires_at").Comment("过期时间"),
		field.String("signed_url").Optional().Comment("已签约短链接"),
		field.String("unsigned_url").Optional().Comment("已签约短链接"),
		field.JSON("paths", &model.Paths{}).Comment("文档各项路径"),
		field.Time("create_at").Comment("创建时间"),
	}
}

// Edges of the Document.
func (Document) Edges() []ent.Edge {
	return nil
}

func (Document) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("hash"),
		index.Fields("template_id"),
		index.Fields("id_card_number"),
		index.Fields("expires_at"),
	}
}
