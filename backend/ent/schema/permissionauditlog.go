package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"entgo.io/ent/schema/mixin"
)

// PermissionAuditLog records security-sensitive RBAC changes.
type PermissionAuditLog struct{ ent.Schema }

func (PermissionAuditLog) Mixin() []ent.Mixin { return []ent.Mixin{mixin.CreateTime{}} }
func (PermissionAuditLog) Fields() []ent.Field {
	return []ent.Field{field.Int("operator_id"), field.String("action").MaxLen(40), field.String("target_type").MaxLen(40), field.Int("target_id"), field.String("request_id").Default("").MaxLen(100), field.String("ip_address").Default("").MaxLen(64), field.JSON("details", map[string]any{}).Optional()}
}
func (PermissionAuditLog) Indexes() []ent.Index {
	return []ent.Index{index.Fields("operator_id", "create_time"), index.Fields("target_type", "target_id")}
}
