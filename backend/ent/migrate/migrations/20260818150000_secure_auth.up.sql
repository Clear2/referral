ALTER TABLE "users" ADD COLUMN "password_hash" character varying NULL;

INSERT INTO "permissions" ("create_time", "update_time", "name", "code", "module", "description")
VALUES
  (CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, '查看推荐管理', 'referral:read', 'referral', '查看推荐关系、积分流水和汇总'),
  (CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, '创建推荐配置', 'referral:create', 'referral', '创建和发放推荐配置'),
  (CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, '修改推荐配置', 'referral:update', 'referral', '修改推荐配置'),
  (CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, '删除推荐配置', 'referral:delete', 'referral', '删除推荐配置')
ON CONFLICT ("code") DO NOTHING;

INSERT INTO "role_permissions" ("role_id", "permission_id")
SELECT r.id, p.id FROM "roles" r CROSS JOIN "permissions" p WHERE r.code = 'super_admin'
ON CONFLICT DO NOTHING;
