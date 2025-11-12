--  Copyright (c) 2025 Michael D Henderson. All rights reserved.

-- Add the "user" role to the roles table
INSERT INTO roles (role_id, is_active, description, created_at, updated_at)
VALUES ('user', 1, 'regular user role', 0, 0);

-- Assign "user" role to all users who have the "active" role (excluding sysop)
INSERT INTO user_roles (user_id, role_id, created_at, updated_at)
SELECT ur.user_id, 'user', ur.created_at, ur.updated_at
FROM user_roles ur
WHERE ur.role_id = 'active'
  AND NOT EXISTS (SELECT 1
                  FROM user_roles ur2
                  WHERE ur2.user_id = ur.user_id
                    AND ur2.role_id = 'user')
  AND NOT EXISTS (SELECT 1
                  FROM user_roles ur3
                  WHERE ur3.user_id = ur.user_id
                    AND ur3.role_id = 'sysop');

-- Assign "guest" role to all users who do NOT have the "active" role (excluding sysop)
INSERT INTO user_roles (user_id, role_id, created_at, updated_at)
SELECT u.user_id, 'guest', u.created_at, u.updated_at
FROM users u
WHERE NOT EXISTS (SELECT 1
                  FROM user_roles ur
                  WHERE ur.user_id = u.user_id
                    AND ur.role_id = 'active')
  AND NOT EXISTS (SELECT 1
                  FROM user_roles ur2
                  WHERE ur2.user_id = u.user_id
                    AND ur2.role_id = 'guest')
  AND NOT EXISTS (SELECT 1
                  FROM user_roles ur3
                  WHERE ur3.user_id = u.user_id
                    AND ur3.role_id = 'sysop');
