-- dashboard perm view
CREATE
OR REPLACE VIEW dashboard_perms AS
SELECT
    urd.user_id as user_id,
    d.id as dash_id,
    d.name as dash_name,
    d.description as dash_description,
    r.id as role_id,
    r.name as role_name,
    p.id as perm_id,
    p.name as perm_name
FROM
    roles r,
    permissions p,
    role_has_permissions rhp,
    user_role_dashboard urd,
    dashboard d
WHERE
    urd.role_id = rhp.role_id
    AND rhp.role_id = r.id
    AND rhp.permission_id = p.id
    AND urd.dashboard_id = d.id;

-- view perm view
CREATE
OR REPLACE VIEW view_perms AS
SELECT
    v.dashboard_id as dash_id,
    urv.user_id as user_id,
    v.id as view_id,
    v.name as view_name,
    v.description as view_desc,
    r.id as role_id,
    r.name as role_name,
    p.id as perm_id,
    p.name as perm_name
FROM
    roles r,
    permissions p,
    role_has_permissions rhp,
    user_role_view urv,
    view v
WHERE
    urv.role_id = rhp.role_id
    AND rhp.role_id = r.id
    AND rhp.permission_id = p.id
    AND urv.view_id = v.id;