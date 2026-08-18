package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"entgo.io/ent/schema/mixin"
)

type CreditTransaction struct{ ent.Schema }

func (CreditTransaction) Mixin() []ent.Mixin { return []ent.Mixin{mixin.Time{}} }

func (CreditTransaction) Fields() []ent.Field {
	return []ent.Field{
		field.Int("user_id"),
		field.Int("referral_id").Unique(),
		field.Int("amount"),
		field.String("reason"),
	}
}

func (CreditTransaction) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).Ref("credit_transactions").Field("user_id").Unique().Required(),
		edge.From("referral", Referral.Type).Ref("credit_transaction").Field("referral_id").Unique().Required(),
	}
}

func (CreditTransaction) Indexes() []ent.Index { return []ent.Index{index.Fields("user_id")} }
