# Ember Data Integration

## Overview

This document outlines the plan for integrating Ember Data into the OttoApp frontend to consume JSON:API responses from the backend.

## Constraints and Guidelines

1. **Do not break Ember Simple Auth** - ESA handles authentication via cookies; Ember Data handles data fetching. These are separate concerns and should not interfere with each other.

2. **Start with `/user/profile`** - Focus on the `/api/my/profile` route as the first implementation.

3. **Use GJS only** - No HBS templates or components. Stick with GJS for consistency with Ember v6.8.

4. **Make small, incremental changes** - Pause frequently for review and questions.

## How Ember Data Will Work

### Current Flow (without Ember Data)

1. Route calls `this.api.getProfile()` (custom service)
2. Returns plain JavaScript object
3. Template passes to component as `@model`

### With Ember Data

1. Define a `user` model (class with attributes matching JSON:API response)
2. Install Ember Data and JSON:API adapter (already speaks JSON:API)
3. Route calls `this.store.findRecord('user', 'me')` or a custom query
4. Ember Data fetches from `/api/my/profile`, deserializes JSON:API, returns model instance
5. Template uses the same `@model`, but now it's a tracked Ember Data model
6. Changes auto-sync, caching works automatically

## Implementation Steps

### Step 1: Install Ember Data packages ✅

**Decision: Use existing packages at version 5.6.0**

The frontend already has the required packages installed:
- `ember-data@~5.6.0`
- `@warp-drive/ember@~5.6.0`

These packages provide:
- **ember-data** - Core library (models, adapters, serializers, store)
- **@warp-drive/ember** - Ember-specific utilities and async data components
- **Built-in JSON:API support** - No additional packages needed

**Why not upgrade to 5.7.0?**
- Ember Data versioning is complex; 5.7.0 versions exist only for `@warp-drive/*` packages
- The main `ember-data` package skipped from 5.3 LTS to later versions
- 5.6.0 is stable and contains all features needed for JSON:API integration
- No breaking changes or critical features in 5.7.0 that affect our use case

**Status:** Complete - proceeding with existing 5.6.0 packages

### Step 2: Create User Model ✅

**Status:** Complete - model generated using `godel` tool

Created `frontend/app/models/user.js` with all attributes matching backend JSON:API response:
- `username` (string)
- `email` (string)
- `timezone` (string)
- `roles` (array)
- `permissions` (object/map)
- `createdAt` (date)
- `updatedAt` (date)

**Tool:** `cmd/godel` CLI generates models from Go structs with `jsonapi` tags
- Keeps backend and frontend in sync
- Prevents schema drift
- See commit da89734

### Step 3: Configure Adapter/Serializer ✅

**Status:** Complete - backend is fully JSON:API compliant

Created `frontend/app/adapters/application.js`:
- Uses `JSONAPIAdapter` from `@ember-data/adapter/json-api`
- Sets `namespace = 'api'` for all requests
- Adds CSRF token via `X-CSRF-Token` header from session
- Works seamlessly with Ember Simple Auth

**Backend compatibility verified:**
- ✅ Returns `Content-Type: application/vnd.api+json`
- ✅ Proper JSON:API structure with `data`, `type`, `id`, `attributes`
- ✅ ISO8601 dates, arrays, and objects properly formatted
- ✅ Error responses follow JSON:API error spec

**No serializer needed** - backend uses `hashicorp/jsonapi` library which outputs spec-compliant responses

**Type naming convention:**
- Backend uses **singular** type: `type: "user"` (not `"users"`)
- This deviates from JSON:API convention (which recommends plural)
- Decision made because Ember Data's inflector not auto-loading with Vite/Embroider
- Simpler: singular everywhere (model file, class name, type name)

### Step 4: Update Route Model Hook ✅

**Status:** Complete - route and controller use Ember Data

Updated `frontend/app/routes/user/profile.js`:
- Changed from `@service api` to `@service store`
- Replaced `this.api.getProfile()` with `this.store.findRecord('user', userId)`
- Gets userId from `session.data.authenticated.user.id`

Updated `frontend/app/controllers/user/profile.js`:
- Removed `@service api` dependency
- Direct model attribute updates: `this.model.email = changes.email`
- Call `await this.model.save()` - automatic PATCH request
- Automatic rollback on error: `this.model.rollbackAttributes()`

**Benefits:**
- Automatic dirty tracking - only changed attributes sent to server
- CSRF token automatically added by adapter
- Response automatically updates the model
- Cleaner code with less boilerplate

### Step 5: Test ✅

**Status:** Complete - fully functional

**Verified:**
- ✅ Profile page loads with Ember Data
- ✅ Update works (sends only changed attributes via PATCH)
- ✅ Ember Simple Auth integration intact
- ✅ Session handling unchanged
- ✅ CSRF tokens working
- ✅ Error handling and rollback working
- ✅ Success messages displaying correctly

**Backend adjustments:**
- Changed to `type: "users"` (plural) per JSON:API convention
- Changed timestamps to dasherized: `created-at`, `updated-at`
- Updated `UserPatchPayload` to match
- Added `godel` support for dasherized-to-camelCase conversion

## Benefits

- **Automatic caching** - No manual cache management
- **Automatic change tracking** - Models are reactive
- **JSON:API compliance** - Native support for spec
- **Relationship handling** - Future-proof for related data
- **Less boilerplate** - No custom API service methods needed
