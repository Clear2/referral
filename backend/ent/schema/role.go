package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"entgo.io/ent/schema/mixin"
)

// Role groups permissions that can be assigned to users.
type Role struct{ ent.Schema }

func (Role) Mixin() []ent.Mixin { return []ent.Mixin{mixin.Time{}} }

func (Role) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").NotEmpty().MaxLen(80),
		field.String("code").NotEmpty().MaxLen(80).Unique(),
		field.String("description").Default("").MaxLen(500),
		field.Bool("enabled").Default(true),
		field.Bool("is_system").Default(false),
	}
}

func (Role) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("users", User.Type).Ref("roles"),
		edge.To("permissions", Permission.Type),
		edge.To("menus", Menu.Type),
	}
}

func (Role) Indexes() []ent.Index { return []ent.Index{index.Fields("enabled")} }
