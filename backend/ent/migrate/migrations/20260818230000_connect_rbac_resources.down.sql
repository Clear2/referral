DELETE FROM "role_menus" rm
USING "menus" m
WHERE rm.menu_id = m.id AND m.component IN ('routes/dashboard', 'routes/users', 'routes/administrators', 'routes/permissions');

DELETE FROM "permission_menus" pm
USING "menus" m
WHERE pm.menu_id = m.id AND m.component IN ('routes/dashboard', 'routes/users', 'routes/administrators', 'routes/permissions');

DELETE FROM "menus" WHERE component IN ('routes/dashboard', 'routes/users', 'routes/administrators', 'routes/permissions');
