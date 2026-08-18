package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"entgo.io/ent/schema/mixin"
)

// PermissionGroup organizes permission points by business capability.
type PermissionGroup struct{ ent.Schema }

func (PermissionGroup) Mixin() []ent.Mixin { return []ent.Mixin{mixin.Time{}} }
func (PermissionGroup) Fields() []ent.Field {
	return []ent.Field{field.String("name").NotEmpty().MaxLen(100), field.String("module").NotEmpty().MaxLen(60), field.String("description").Default("").MaxLen(500), field.Int("parent_id").Optional().Nillable(), field.Int("sort_order").Default(0), field.Bool("enabled").Default(true)}
}
func (PermissionGroup) Edges() []ent.Edge {
	return []ent.Edge{edge.From("parent", PermissionGroup.Type).Ref("children").Field("parent_id").Unique(), edge.To("children", PermissionGroup.Type), edge.To("permissions", Permission.Type)}
}
func (PermissionGroup) Indexes() []ent.Index {
	return []ent.Index{index.Fields("parent_id"), index.Fields("module")}
}
