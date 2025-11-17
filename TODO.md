# TODO

## Sprint 13

### Backend

[x] Create JSON:API compliant version service

### Frontend

[x] Add Settings page with version info
  [x] Check if Settings link exists in sidebar navigation
    [x] File: frontend/app/components/layouts/home-screens/sidebar/sidebar.gjs
    [x] Add Settings menu item if missing
  [x] Create admin settings routes (nested)
    [x] File: frontend/app/routes/admin/settings/index.js - redirects to account
    [x] File: frontend/app/routes/admin/settings/account.js - loads user model
    [x] File: frontend/app/routes/admin/settings/notifications.js - placeholder
    [x] File: frontend/app/routes/admin/settings/about.js - loads version model
  [x] Create user settings routes (nested)
    [x] File: frontend/app/routes/user/settings/index.js - redirects to account
    [x] File: frontend/app/routes/user/settings/account.js - loads user model
    [x] File: frontend/app/routes/user/settings/maps.js - placeholder
    [x] File: frontend/app/routes/user/settings/teams.js - placeholder
    [x] File: frontend/app/routes/user/settings/notifications.js - placeholder
    [x] File: frontend/app/routes/user/settings/about.js - loads version model
  [x] Create admin settings template
    [x] File: frontend/app/templates/admin/settings.gjs - parent with tab nav
    [x] File: frontend/app/templates/admin/settings/account.gjs
    [x] File: frontend/app/templates/admin/settings/notifications.gjs
    [x] File: frontend/app/templates/admin/settings/about.gjs
    [x] Tabs: Account | Notifications | About
  [x] Create user settings template
    [x] File: frontend/app/templates/user/settings.gjs - parent with tab nav
    [x] File: frontend/app/templates/user/settings/account.gjs
    [x] File: frontend/app/templates/user/settings/maps.gjs
    [x] File: frontend/app/templates/user/settings/teams.gjs
    [x] File: frontend/app/templates/user/settings/notifications.gjs
    [x] File: frontend/app/templates/user/settings/about.gjs
    [x] Tabs: Account | Maps | Teams | Notifications | About
  [x] Create shared setting components (used by both admin and user)
    [x] File: frontend/app/components/settings/account.gjs - under construction
    [x] File: frontend/app/components/settings/notifications.gjs - under construction
    [x] File: frontend/app/components/settings/about.gjs - displays version info
    [x] Display @version.full and @version.short
    [x] Use grid section layout matching Tailwind UI pattern
  [x] Create user-specific setting components
    [x] File: frontend/app/components/settings/maps.gjs - under construction
    [x] File: frontend/app/components/settings/teams.gjs - under construction
  [x] Create under-construction component
    [x] File: frontend/app/components/settings/under-construction.gjs

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
