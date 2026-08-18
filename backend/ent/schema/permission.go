package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"entgo.io/ent/schema/mixin"
)

// Permission is an atomic operation granted through a role.
type Permission struct{ ent.Schema }

func (Permission) Mixin() []ent.Mixin { return []ent.Mixin{mixin.Time{}} }

func (Permission) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").NotEmpty().MaxLen(100),
		field.String("code").NotEmpty().MaxLen(120).Unique(),
		field.String("module").Default("system").MaxLen(60),
		field.String("description").Default("").MaxLen(500),
		field.Bool("enabled").Default(true),
		field.Int("group_id").Optional().Nillable(),
	}
}

func (Permission) Edges() []ent.Edge {
	return []ent.Edge{edge.From("roles", Role.Type).Ref("permissions"), edge.From("group", PermissionGroup.Type).Ref("permissions").Field("group_id").Unique(), edge.To("apis", APIResource.Type), edge.To("menus", Menu.Type)}
}

func (Permission) Indexes() []ent.Index {
	return []ent.Index{index.Fields("module"), index.Fields("enabled")}
}
