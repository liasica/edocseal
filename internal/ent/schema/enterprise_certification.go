package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// EnterpriseCertification holds the schema definition for the EnterpriseCertification entity.
type EnterpriseCertification struct {
	ent.Schema
}

// Annotations of the EnterpriseCertification.
func (EnterpriseCertification) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "enterprise_certification"},
		entsql.WithComments(true),
	}
}

// Fields of the EnterpriseCertification.
func (EnterpriseCertification) Fields() []ent.Field {
	return []ent.Field{
		field.String("credit_code").MaxLen(18).Comment("统一社会信用代码"),
		field.String("issuer").MaxLen(32).Comment("证书生成方式"),
		field.String("private_path").MaxLen(255).Comment("私钥路径"),
		field.String("cert_path").MaxLen(255).Comment("证书路径"),
		field.Time("expires_at").Comment("证书过期时间"),
	}
}

// Edges of the EnterpriseCertification.
func (EnterpriseCertification) Edges() []ent.Edge {
	return nil
}

// Indexes of the EnterpriseCertification.
func (EnterpriseCertification) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("credit_code", "issuer").Unique(),
	}
}
