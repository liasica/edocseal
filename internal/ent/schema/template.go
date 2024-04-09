package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
)

// Template holds the schema definition for the Template entity.
type Template struct {
	ent.Schema
}

// Annotations of the Template.
func (Template) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "template"},
		entsql.WithComments(true),
	}
}

// Fields of the Template.
func (Template) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").NotEmpty().Comment("模板ID"),
		field.String("name").NotEmpty().Comment("模板名称"),
		field.String("file").NotEmpty().Comment("模板文件"),
		field.JSON("fields", map[string][4]float64{}).Comment("模板字段，从左往右从上往下坐标"),
		field.Time("created_at").Immutable().Default(time.Now).Comment("创建时间"),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now).Comment("更新时间"),
	}
}

// Edges of the Template.
func (Template) Edges() []ent.Edge {
	return nil
}
