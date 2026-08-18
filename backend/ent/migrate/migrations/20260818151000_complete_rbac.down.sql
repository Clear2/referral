DROP TABLE "permission_audit_logs";
DROP TABLE "permission_menus";
DROP TABLE "permission_apis";
DROP TABLE "api_resources";
ALTER TABLE "menus" DROP COLUMN "redirect", DROP COLUMN "component";
ALTER TABLE "permissions" DROP COLUMN "group_id";
DROP TABLE "permission_groups";
