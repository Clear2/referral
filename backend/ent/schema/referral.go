package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"entgo.io/ent/schema/mixin"
)

type Referral struct{ ent.Schema }

func (Referral) Mixin() []ent.Mixin { return []ent.Mixin{mixin.Time{}} }

func (Referral) Fields() []ent.Field {
	return []ent.Field{
		field.Int("inviter_id"),
		field.Int("invitee_id").Unique(),
		field.Int("reward").Positive(),
	}
}

func (Referral) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("inviter", User.Type).Ref("sent_referrals").Field("inviter_id").Unique().Required(),
		edge.From("invitee", User.Type).Ref("received_referral").Field("invitee_id").Unique().Required(),
		edge.To("credit_transaction", CreditTransaction.Type).Unique(),
	}
}

func (Referral) Indexes() []ent.Index { return []ent.Index{index.Fields("inviter_id")} }
