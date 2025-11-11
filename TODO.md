# TODO

## Sprint 7

Implement two authenticated roles: admin and user with RBAC.

### Backend

[ ] Create migration script to add "user" role to roles table and assign roles to existing users:
  [ ] Add "user" role to roles table
  [ ] Assign "user" role to all users with "active" role
  [ ] Assign "guest" role to any users without "active" role
[ ] Update expectedSchemaVersion in backend/stores/sqlite/db.go to match new migration
[ ] Add role field to user creation (default to "user")
[ ] Add command to update user role via CLI
[ ] Implement profile management endpoints:
  [ ] `GET /api/users/me` - Get current user's profile
  [ ] `GET /api/users/:id` - Get user profile (with authorization check)
  [ ] `PATCH /api/users/:id` - Update profile (enforce username edit restrictions)
  [ ] `PUT /api/users/:id/password` - Update own password (requires current password)
  [ ] `POST /api/users/:id/reset-password` - Admin resets user password
[ ] Implement admin endpoints:
  [ ] `GET /api/users` - List all non-admin users (admin only)
  [ ] `POST /api/users` - Create user with role (admin only)
  [ ] `PATCH /api/users/:id/role` - Update user role (admin only)
[ ] Add RBAC authorization logic:
  [ ] Regular users can edit own profile (except username)
  [ ] Regular users can update own password
  [ ] Admins can edit any non-admin user's profile (including username)
  [ ] Admins can reset password for any non-admin user
[ ] Return 403 for unauthorized actions, 404 for inaccessible resources
[ ] Include permissions metadata in user profile responses

### Frontend

[ ] Create new routes:
  [ ] `/users` - User dashboard (requires "user" role)
  [ ] `/users/profile` - Edit own profile
  [ ] `/users/password` - Change own password
  [ ] `/admin` - Admin dashboard (requires "admin" role)
  [ ] `/admin/users` - List all users
  [ ] `/admin/users/:id` - View/edit user profile
  [ ] `/admin/users/:id/reset-password` - Reset user password
  [ ] `/admin/users/new` - Create new user
[ ] Add authorization service or extend session service:
  [ ] `hasRole(role)` - Check if current user has role
  [ ] `canEditUsername()` - Permission helper
  [ ] `canAccessAdminRoutes()` - Role check helper
[ ] Implement route protection with role checks
[ ] Create components:
  [ ] User profile form (conditionally read-only username)
  [ ] Password change form (requires current password)
  [ ] Admin user list table
  [ ] Admin user form (with role selector)
  [ ] Admin password reset dialog
[ ] Update navigation to show/hide links based on roles
[ ] Add API service methods for new endpoints


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
