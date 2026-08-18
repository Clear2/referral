DELETE FROM "role_permissions" WHERE "permission_id" IN (SELECT "id" FROM "permissions" WHERE "code" IN ('referral:read', 'referral:create', 'referral:update', 'referral:delete'));
DELETE FROM "permissions" WHERE "code" IN ('referral:read', 'referral:create', 'referral:update', 'referral:delete');
ALTER TABLE "users" DROP COLUMN "password_hash";
