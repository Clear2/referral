INSERT INTO "menus" ("create_time", "update_time", "name", "path", "icon", "component", "type", "sort_order", "enabled")
SELECT CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, v.name, v.path, v.icon, v.component, 'MENU', v.sort_order, true
FROM (VALUES
  ('邀请概览', '/', 'LayoutDashboard', 'routes/dashboard', 10),
  ('普通用户', '/users', 'Users', 'routes/users', 20),
  ('管理账号', '/administrators', 'UserCog', 'routes/administrators', 30),
  ('权限管理', '/permissions', 'ShieldCheck', 'routes/permissions', 40)
) AS v(name, path, icon, component, sort_order)
WHERE NOT EXISTS (SELECT 1 FROM "menus" m WHERE m.path = v.path);

INSERT INTO "permission_menus" ("permission_id", "menu_id")
SELECT p.id, m.id
FROM "permissions" p
JOIN "menus" m ON
  (p.code = 'referral:read' AND m.path = '/') OR
  (p.code = 'system:rbac' AND m.path IN ('/users', '/administrators', '/permissions'))
ON CONFLICT DO NOTHING;

INSERT INTO "role_menus" ("role_id", "menu_id")
SELECT DISTINCT rp.role_id, pm.menu_id
FROM "role_permissions" rp
JOIN "permission_menus" pm ON pm.permission_id = rp.permission_id
ON CONFLICT DO NOTHING;
