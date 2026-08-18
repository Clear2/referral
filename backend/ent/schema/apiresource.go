package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"entgo.io/ent/schema/mixin"
)

// APIResource describes one HTTP operation protected by permissions.
type APIResource struct{ ent.Schema }

func (APIResource) Mixin() []ent.Mixin { return []ent.Mixin{mixin.Time{}} }
func (APIResource) Fields() []ent.Field {
	return []ent.Field{field.String("name").NotEmpty().MaxLen(160), field.String("method").NotEmpty().MaxLen(10), field.String("path").NotEmpty().MaxLen(240), field.String("description").Default("").MaxLen(500), field.Bool("enabled").Default(true)}
}
func (APIResource) Edges() []ent.Edge {
	return []ent.Edge{edge.From("permissions", Permission.Type).Ref("apis")}
}
func (APIResource) Indexes() []ent.Index {
	return []ent.Index{index.Fields("method", "path").Unique(), index.Fields("enabled")}
}
