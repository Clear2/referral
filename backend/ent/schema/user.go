package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"entgo.io/ent/schema/mixin"
)

// User holds the schema definition for the User entity.
type User struct {
	ent.Schema
}

// Annotations of the User.
func (User) Annotations() []schema.Annotation {
	return []schema.Annotation{
		// Adding this annotation to the schema enables
		// comments for the table and all its fields.
		entsql.WithComments(true),
		schema.Comment("User table stores all user information."),
	}
}

// Fields of the User.
func (User) Fields() []ent.Field {
	return []ent.Field{
		field.String("name"),
		field.String("email").
			Unique(),
		field.String("referral_code").Optional().Nillable().Unique(),
		field.Int("credit_balance").Default(0).NonNegative(),
		field.Bool("enabled").Default(true),
		field.String("password_hash").Optional().Sensitive(),
	}
}

// Mixin of the User.
func (User) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.Time{},
	}
}

// Edges of the User.
func (User) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("sent_referrals", Referral.Type),
		edge.To("received_referral", Referral.Type).Unique(),
		edge.To("credit_transactions", CreditTransaction.Type),
		edge.To("roles", Role.Type),
	}
}

// Indexes of the User.
func (User) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("name", "email"),
	}
}
