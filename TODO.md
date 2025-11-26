# TODO

## Sprint 17

### GM Turn Report Uploads
- [ ] **Backend:** Verify/Update `gm` role in `backend/auth/policies.go` (CanUploadTurnReports).
- [ ] **Backend:** Implement `PostDocument` in `backend/servers/rest/documents.go`.
    - [ ] Enforce `gm` role (or `admin`/`sysop`).
    - [ ] Enforce max file size (150KB).
    - [ ] Parse clan heading using `documentsSvc.ParseClanHeading`.
    - [ ] Store document metadata and contents in SQLite.
    - [ ] Return JSON:API response (201 Created).
- [ ] **Frontend:** Install `ember-concurrency`.
- [ ] **Frontend:** Create `user.uploads` route.
- [ ] **Frontend:** Create `<TurnReportUploader />` component (`.gjs`).
    - [ ] Support multiple file selection/drop.
    - [ ] Implement sequential upload queue using `ember-concurrency`.
    - [ ] Display status for each file.
- [ ] **Testing:** Verify with `penguin` (GM) and `catbird` (User) using `./tmp/gemini` DB.
