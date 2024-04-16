package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
)

// Certification holds the schema definition for the Certification entity.
type Certification struct {
	ent.Schema
}

// Annotations of the Certification.
func (Certification) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "certification"},
		entsql.WithComments(true),
	}
}

// Fields of the Certification.
func (Certification) Fields() []ent.Field {
	return []ent.Field{
		field.String("id_card_number").Unique().MaxLen(18).Comment("身份证号码"),
		field.String("private_path").MaxLen(255).Comment("私钥路径"),
		field.String("cert_path").MaxLen(255).Comment("证书路径"),
		field.Time("expires_at").Comment("证书过期时间"),
	}
}

// Edges of the Certification.
func (Certification) Edges() []ent.Edge {
	return nil
}
