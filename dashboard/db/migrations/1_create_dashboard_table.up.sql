DROP TABLE IF EXISTS dashboard;

-- Dashboard and views inside dashboard
CREATE TABLE dashboard(
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    name text NOT NULL,
    description text NOT NULL,
    created_at timestamp NOT NULL DEFAULT NOW()
);

CREATE TABLE view (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    dashboard_id uuid NOT NULL REFERENCES dashboard(id) ON DELETE CASCADE,
    name text NOT NULL,
    description text NOT NULL,
    created_at timestamp NOT NULL DEFAULT NOW()
);

CREATE TABLE roles(
    id SERIAL PRIMARY KEY,
    name text NOT NULL UNIQUE
);

INSERT INTO
    roles (name)
VALUES
    ('admin'),
    ('editor'),
    ('viewer'),
    ('commenter'),
    ('manager');

-- permissions and roles
CREATE TABLE permissions(
    id SERIAL PRIMARY KEY,
    name text NOT NULL UNIQUE
);

INSERT INTO
    permissions (name)
VALUES
    ('edit_access'),
    ('delete'),
    ('edit'),
    ('comment'),
    ('read');

CREATE TABLE role_has_permissions(
    role_id int NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    permission_id int NOT NULL REFERENCES permissions(id) ON DELETE CASCADE,
    PRIMARY KEY (role_id, permission_id)
);

-- admin has all permissions
INSERT INTO
    role_has_permissions (role_id, permission_id)
VALUES
    (1, 1),
    (1, 2),
    (1, 3),
    (1, 4),
    (1, 5);

-- editor has all content based permissions except delete
INSERT INTO
    role_has_permissions (role_id, permission_id)
VALUES
    (2, 3),
    (2, 4),
    (2, 5);

-- viewer has read permission
INSERT INTO
    role_has_permissions (role_id, permission_id)
VALUES
    (3, 5);

-- commenter has comment and read permission
INSERT INTO
    role_has_permissions (role_id, permission_id)
VALUES
    (4, 4),
    (4, 5);

-- manager has edit access and read permission
INSERT INTO
    role_has_permissions (role_id, permission_id)
VALUES
    (5, 1),
    (5, 2),
    (5, 3),
    (5, 4),
    (5, 5);

-- relating users to roles for dashboard
CREATE TABLE user_role_dashboard(
    user_id uuid NOT NULL,
    role_id int NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    dashboard_id uuid NOT NULL REFERENCES dashboard(id) ON DELETE CASCADE,
    PRIMARY KEY (user_id, dashboard_id) -- user can have only one role per dashboard
);

-- relating users to roles for view
CREATE TABLE user_role_view(
    user_id uuid NOT NULL,
    role_id int NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    view_id uuid NOT NULL REFERENCES view(id) ON DELETE CASCADE,
    PRIMARY KEY (user_id, view_id) -- user can have only one role per view
);