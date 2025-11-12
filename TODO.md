# TODO

## Sprint 7

Implement two authenticated roles: admin and user with RBAC.

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

[ ] Refactor frontend to use DDAU (Data Down, Actions Up) architecture with Ember v6.8:
    [x] Phase 1: Enhance API Service (app/services/api.js):
        [x] Add getProfile() method - GET /api/profile
        [x] Add updateProfile(data) method - POST /api/profile
        [x] Add getTimezones() method - GET /api/timezones
    [ ] Phase 2: Routes (app/routes/*/*.js) - Data loading only:
        [ ] Fix app/routes/user/profile.js to use API service instead of fetch()
        [ ] Verify app/routes/users/profile.js uses API service
        [ ] Verify app/routes/users/password.js uses API service
        [ ] Verify all admin routes use API service
    [ ] Phase 3: Controllers (app/controllers/*/*.js) - Actions and state:
        [ ] Create app/controllers/user/profile.js (currently missing)
        [ ] Verify app/controllers/users/profile.js follows DDAU
        [ ] Verify app/controllers/users/password.js follows DDAU
        [ ] Verify all admin controllers follow DDAU
    [ ] Phase 4: Templates (app/templates/*/*.gjs) - Import components, pass args:
        [ ] Update app/templates/user/profile.gjs to pass controller actions
        [ ] Update app/templates/users/profile.gjs to pass all needed args
        [ ] Create app/templates/users/password.gjs if missing
        [ ] Update all admin templates to pass proper args
    [ ] Phase 5: Components (app/components/**/*.gjs) - Pure presentation:
        [ ] Refactor app/components/timezone-picker.gjs to remove fetch(), receive @timezones
        [ ] Refactor app/components/user/profile.gjs to remove fetch/save logic
        [ ] Refactor app/components/users/profile.gjs to be pure presentational
        [ ] Refactor app/components/users/password.gjs to be pure presentational
        [ ] Refactor app/components/admin/users/index.gjs to be pure presentational
        [ ] Refactor app/components/admin/users/edit.gjs to be pure presentational
        [ ] Refactor app/components/admin/users/new.gjs to be pure presentational
        [ ] Verify only api.js and ESA authenticator call fetch()

### Stretch Goals

[ ] Fix API runner test for /api/users/me endpoint (sessionsSvc integration issue)
[ ] Complete API endpoint testing with live server

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
