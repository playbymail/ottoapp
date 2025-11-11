# TODO

## Sprint 7

Implement two authenticated roles: admin and user with RBAC.

### Stretch Goals

[ ] Fix API runner test for /api/users/me endpoint (sessionsSvc integration issue)
[ ] Complete API endpoint testing with live server

## Completed Sprints

### Sprint 7 - RBAC Implementation

### Backend

[x] Create migration script to add "user" role to roles table and assign roles to existing users:
  [x] Add "user" role to roles table
  [x] Assign "user" role to all users with "active" role
  [x] Assign "guest" role to any users without "active" role
[x] Update expectedSchemaVersion in backend/stores/sqlite/db.go to match new migration
[x] Add role field to user creation (default to "user")
[x] Add command to update user role via CLI
[x] Implement profile management endpoints:
  [x] `GET /api/users/me` - Get current user's profile
  [x] `GET /api/users/:id` - Get user profile (with authorization check)
  [x] `PATCH /api/users/:id` - Update profile (enforce username edit restrictions)
  [x] `PUT /api/users/:id/password` - Update own password (requires current password)
  [x] `POST /api/users/:id/reset-password` - Admin resets user password
[x] Implement admin endpoints:
  [x] `GET /api/users` - List all non-admin users (admin only)
  [x] `POST /api/users` - Create user with role (admin only)
  [x] `PATCH /api/users/:id/role` - Update user role (admin only)
[x] Add RBAC authorization logic:
  [x] Regular users can edit own profile (except username)
  [x] Regular users can update own password
  [x] Admins can edit any non-admin user's profile (including username)
  [x] Admins can reset password for any non-admin user
[x] Return 403 for unauthorized actions, 404 for inaccessible resources
[x] Include permissions metadata in user profile responses

### Frontend

[x] Create new routes:
  [x] `/users` - User dashboard (requires "user" role)
  [x] `/users/profile` - Edit own profile
  [x] `/users/password` - Change own password
  [x] `/admin` - Admin dashboard (requires "admin" role)
  [x] `/admin/users` - List all users
  [x] `/admin/users/:id` - View/edit user profile
  [x] `/admin/users/:id/reset-password` - Reset user password
  [x] `/admin/users/new` - Create new user
[x] Add authorization service or extend session service:
  [x] `hasRole(role)` - Check if current user has role
  [x] `canEditUsername()` - Permission helper
  [x] `canAccessAdminRoutes()` - Role check helper
[x] Implement route protection with role checks
[x] Create components:
  [x] User profile form (conditionally read-only username)
  [x] Password change form (requires current password)
  [x] Admin user list table
  [x] Admin user form (with role selector)
  [x] Admin password reset dialog
[x] Update navigation to show/hide links based on roles
[x] Add API service methods for new endpoints


## Future Sprints

### Backend

[ ] Implement "app" commands to load users, reports, etc. for in-memory database configuration
  [ ] Currently we can start the server using an in-memory database but we can't configure it yet

### Command Line

[ ] Add helper function for common database operations in cmd/ottoapp/main.go
  [ ] Extract pattern of getting `--db` flag and opening database connection
  [ ] Many commands repeat this code: `cmd.Flags().GetString("db")` followed by `sqlite.Open()`
  [ ] Helper should handle errors consistently and reduce boilerplate

### Frontend

[ ] Add footer component showing version information
  [ ] Display backend version (from /api/version or similar)
  [ ] Display frontend version (from package.json or build metadata)
  [ ] Include link to current commit on GitHub repository
  [ ] Should be visible on all pages (place in application template)
