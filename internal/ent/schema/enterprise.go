package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
)

// Enterprise holds the schema definition for the Enterprise entity.
type Enterprise struct {
	ent.Schema
}

// Annotations of the Enterprise.
func (Enterprise) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "enterprise"},
		entsql.WithComments(true),
	}
}

// Fields of the Enterprise.
func (Enterprise) Fields() []ent.Field {
	return []ent.Field{
		field.String("credit_code").Unique().MaxLen(18).Comment("统一社会信用代码"),
		field.String("name").MaxLen(100).Comment("企业名称"),
		field.String("province").MaxLen(20).Comment("省份"),
		field.String("city").MaxLen(20).Comment("城市"),
		field.String("person_name").Optional().MaxLen(50).Comment("代办人姓名，向陕西CA 申请证书时使用"),
		field.String("phone").Optional().MaxLen(20).Comment("代办人手机号，向陕西CA 申请证书时使用"),
		field.String("idcard").Optional().MaxLen(18).Comment("代办人身份证号，向陕西CA 申请证书时使用"),
		field.Bool("is_default").Default(false).Comment("是否默认企业，签约请求未指定企业时使用"),
	}
}

// Edges of the Enterprise.
func (Enterprise) Edges() []ent.Edge {
	return nil
}
