package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"entgo.io/ent/schema/mixin"
)

// Menu describes a navigable admin resource.
type Menu struct{ ent.Schema }

func (Menu) Mixin() []ent.Mixin { return []ent.Mixin{mixin.Time{}} }

func (Menu) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").NotEmpty().MaxLen(100),
		field.String("path").Default("").MaxLen(240),
		field.String("icon").Default("").MaxLen(100),
		field.String("component").Default("").MaxLen(200),
		field.String("redirect").Default("").MaxLen(240),
		field.Enum("type").Values("CATALOG", "MENU", "BUTTON", "EMBEDDED", "LINK").Default("MENU"),
		field.Int("parent_id").Optional().Nillable(),
		field.Int("sort_order").Default(0),
		field.Bool("enabled").Default(true),
	}
}

func (Menu) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("parent", Menu.Type).Ref("children").Field("parent_id").Unique(),
		edge.To("children", Menu.Type),
		edge.From("roles", Role.Type).Ref("menus"),
		edge.From("permissions", Permission.Type).Ref("menus"),
	}
}

func (Menu) Indexes() []ent.Index {
	return []ent.Index{index.Fields("parent_id"), index.Fields("sort_order")}
}
