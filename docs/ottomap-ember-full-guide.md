# OttoMap Ember 6.8 — Routes, Controllers, and Data Flow Guide

This document summarizes best practices and the complete explanation of how routes, controllers, templates, and
components interact in Ember 6.8, using the OttoMap `/uploads` feature as a running example.

---

## 1. Data and Fetching Best Practices

In Ember 6.8, you’ll find three styles of data access in legacy codebases:

1. Components directly calling `fetch()`.
2. Routes or controllers calling `fetch()` and passing data down.
3. Components using an injected `@service api` to make calls.

### ✅ Modern, Recommended Pattern

- **Only one service** (`app/services/api.js`) should talk to the backend.
- **Routes** and **controllers** should use that service.
- **Components** should only render data or call actions passed from parents.
- **Actions** (like `onSave`, `onDelete`, etc.) should bubble *upward* to routes/controllers, never make `fetch()` calls
  directly.

This is Ember’s “Data Down, Actions Up” (DDAU) pattern.

Ember’s architecture emphasizes **data down, actions up (DDAU)**:

```
Route → Controller → Template → Component → Service
                     ↑                 ↓
                 Actions ↑         API fetch ↓
```

Only two places should make network calls:

- **Ember Simple Auth** (session management)
- **API service** (your own fetch wrapper)

---

## 2. Example File Structure for a Fully Built-Out Page

Below is what a **fully built-out Ember route** looks like in `app/` when you include *everything you might ever
reasonably create* for one page in OttoMap (route + controller + template + components + test + style + service hook).

For the route `/uploads/:upload_id`:

```
app/
├── routes/
│   └── uploads/
│       ├── index.js  →  uploads.index   →  /uploads                               
│       ├── show.js   →  uploads.show    →  /uploads/:upload_id     
│       └── edit.js   →  uploads.edit    →  /uploads/:upload_id/edit
├── controllers/
│   └── uploads/
│       └── show.js
├── templates/
│   └── uploads/
│       ├── show.gjs
│       └── index.gjs
├── components/
│   └── uploads/
│       ├── details.gjs
│       ├── file-info.gjs
│       └── share-dialog.gjs
├── services/
│   └── api.js
└── tests/
    ├── acceptance/
    │   └── uploads-show-test.js
    └── integration/
        └── components/
            └── uploads/
                └── details-test.js
```

---

## 3. Route Lifecycle and Context

When Ember starts to render a route, it performs this sequence:

1. The **Router** matches the URL to a route (e.g. `/uploads/123` → `uploads.show`).
2. Ember instantiates the route object (singleton).
3. It calls, in order:
    - `beforeModel()`
    - `model(params)` → fetch data via API service.
    - `afterModel(model)`
    - `setupController(controller, model)` → assigns `controller.model = model`.
4. The **controller** becomes the **template’s context**.
5. The **template** renders and passes data to **components** via arguments.

### Controller Role

- Holds the route’s model.
- Defines user actions.
- Coordinates navigation or service calls.

### Component Role

- Receives arguments (`@upload`, `@onShare`).
- Maintains local state (`@tracked` vars).
- Calls actions upward.

### Service Role

- Provides reusable methods to communicate with backend (`fetch`, JSON parsing, etc.).

---

## 4. Typical Flow at Runtime

```
Router → Route(model) → Controller(model/actions) → Template → Component → API Service(fetch)
```

| Layer      | Responsibility                  | Fetch?                  |
|------------|---------------------------------|-------------------------|
| Route      | Load data                       | ✅ Yes (via API service) |
| Controller | Hold data and define actions    | ✅ Maybe                 |
| Template   | Layout and bind data            | ❌                       |
| Component  | Render UI and bubble actions up | ❌                       |
| Service    | Perform network calls           | ✅ Yes                   |

---

## 5. File Examples (Simplified)

### Route — `app/routes/uploads/show.js`

```js
import Route from '@ember/routing/route';
import {service} from '@ember/service';

export default class UploadsShowRoute extends Route {
    @service api;

    async model(params) {
        return this.api.getUpload(params.upload_id);
    }
}
```

### Controller — `app/controllers/uploads/show.js`

```js
import Controller from '@ember/controller';
import {service} from '@ember/service';
import {action} from '@ember/object';

export default class UploadsShowController extends Controller {
    @service api;
    @service router;

    @action
    async deleteUpload() {
        await this.api.deleteUpload(this.model.id);
        this.router.transitionTo('uploads.index');
    }

    @action
    async shareUpload(allies) {
        await this.api.shareUpload(this.model.id, allies);
        alert('Shared successfully!');
    }
}
```

### Template — `app/templates/uploads/show.gjs`

```gjs
<template>
  <h1>Upload: {{this.model.filename}}</h1>
  <Uploads::Details @upload={{this.model}} @onShare={{this.shareUpload}} />
</template>
```

### Component — `app/components/uploads/details.gjs`

```gjs
<template>
  <div>
    <Uploads::FileInfo @upload={{@upload}} />
    <button {{on "click" this.openShare}}>Share</button>

    {{#if this.isSharing}}
      <Uploads::ShareDialog
        @upload={{@upload}}
        @onConfirm={{this.confirmShare}}
        @onCancel={{this.cancelShare}}
      />
    {{/if}}
  </div>
</template>

import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { action } from '@ember/object';

export default class UploadsDetailsComponent extends Component {
  @tracked isSharing = false;

  @action openShare() { this.isSharing = true; }
  @action cancelShare() { this.isSharing = false; }

  @action async confirmShare(allies) {
    await this.args.onShare(allies);
    this.isSharing = false;
  }
}
```

### Service — `app/services/api.js`

```js
import Service from '@ember/service';

export default class ApiService extends Service {
    baseUrl = '/api';

    async #request(path, options = {}) {
        const res = await fetch(`${this.baseUrl}${path}`, {
            headers: {'Content-Type': 'application/json'},
            credentials: 'include',
            ...options,
        });
        if (!res.ok) throw new Error(`${res.status} ${res.statusText}`);
        return res.status === 204 ? null : res.json();
    }

    async getUpload(id) {
        return this.#request(`/uploads/${id}`);
    }

    async deleteUpload(id) {
        return this.#request(`/uploads/${id}`, {method: 'DELETE'});
    }

    async shareUpload(id, allies) {
        return this.#request(`/uploads/${id}/share`, {
            method: 'POST',
            body: JSON.stringify({allies}),
        });
    }
}
```

---

## 6. Route Activation Diagram

```
[User visits /uploads/123]
       │
       ▼
Router → Route(model) → API Service(fetch)
       │
       ▼
Controller(model/actions)
       │
       ▼
Template (renders, passes @upload, @onShare)
       │
       ▼
Component (renders, triggers @onShare)
       │
       ▼
Controller action → API Service → backend
```

---

## 7. Route and Outlet Hierarchy

Router map:

```js
this.route('uploads', function () {                 // → app/routes/uploads/index.js
    this.route('show', {path: '/:upload_id'});      // → app/routes/uploads/show.js
    this.route('edit', {path: '/:upload_id/edit'}); // → app/routes/uploads/edit.js
});
```

Ember internally expands this to:

```
application
└── uploads
    ├── uploads.index  → /uploads
    ├── uploads.show   → /uploads/:upload_id
    └── uploads.edit   → /uploads/:upload_id/edit
```

### Template hierarchy (GJS only)

```gjs
// app/templates/application.gjs
<template>
  <ApplicationShell>
    {{outlet}}
  </ApplicationShell>
</template>

// app/templates/uploads.gjs
<template>
  <section class="uploads-layout">
    {{outlet}}
  </section>
</template>

// app/templates/uploads/index.gjs
<template>
  <p>List of uploads goes here…</p>
</template>

// app/templates/uploads/show.gjs
<template>
  <h2>Upload: {{this.model.filename}}</h2>
  <Uploads::Details @upload={{this.model}} />
</template>

// app/templates/uploads/edit.gjs
<template>
  <h2>Edit upload: {{this.model.filename}}</h2>
  <Uploads::EditForm @upload={{this.model}} @onSave={{this.save}} @onCancel={{this.cancel}} />
</template>
```

### File Alignment Diagram

```
application.gjs
└── uploads.gjs          (contains {{outlet}})
    ├── uploads/index.gjs   ← /uploads
    ├── uploads/show.gjs    ← /uploads/123
    └── uploads/edit.gjs    ← /uploads/123/edit
```

---

## 8. TL;DR Summary

- **index** routes are created automatically — they’re the “default child route.”
- **show**, **edit**, etc. are just conventions (no special meaning).
- Routes load data and set up controllers.
- Controllers provide actions and hold models.
- Components render UI and bubble actions up.
- Only the API service and Ember Simple Auth should use `fetch`.

---

**Data Down, Actions Up. Keep `fetch` centralized.**

---

### 💡 Summary of where logic lives

| Concern        | File                                 | Responsibility                                |
|----------------|--------------------------------------|-----------------------------------------------|
| **Route**      | `app/routes/uploads/show.js`         | Fetch model data and set controller state     |
| **Controller** | `app/controllers/uploads/show.js`    | Handle actions that mutate data or navigate   |
| **Template**   | `app/templates/uploads/show.hbs`     | Present layout and wire actions               |
| **Component**  | `app/components/uploads/details.gjs` | Present reusable UI pieces, bubble actions up |
| **Service**    | `app/services/api.js`                | Do all `fetch` I/O                            |
| **Tests**      | `tests/acceptance/...`               | Verify end-to-end behavior                    |
| **Styles**     | `app/styles/uploads.css`             | Page-specific styling                         |

---

## Appendix

Excellent — this is where Ember’s design really shines once you visualize it.
Let’s trace what happens **when a user visits `/uploads/123`** — step by step — showing which file runs and what
data/context each layer receives.

---

## 🔄 Route Activation Flow (simplified)

```
[User visits /uploads/123]
       │
       ▼
🌐 app/router.js
 └── defines routes like this.route('uploads', function () { this.route('show', { path: '/:upload_id' }); });

       │
       ▼
🚦 Ember Router matches → "uploads.show"
       │
       ▼
🧭 Route chain:
   application → uploads → uploads.show
```

---

## 🧠 What Ember does next

### 1️⃣ Route hooks fire (top to bottom)

| Hook                                 | File                                     | Purpose                                                 |
|--------------------------------------|------------------------------------------|---------------------------------------------------------|
| `beforeModel()`                      | `application`, `uploads`, `uploads.show` | Early checks (auth, redirect)                           |
| `model(params)`                      | `uploads.show`                           | Fetch data → calls `api.getUpload(params.upload_id)`    |
| `afterModel(model)`                  | (optional)                               | Run derived logic once data is loaded                   |
| `setupController(controller, model)` | `uploads.show`                           | Assign model to controller (`controller.model = model`) |

➡️ At this point, Ember has the **route model** and the **controller context** ready.

---

### 2️⃣ Controller is prepared

```
app/controllers/uploads/show.js
```

* A controller instance is created (if one exists).
* Ember calls `setupController(controller, model)` → sets `controller.model = model`.
* Any other tracked properties you define are initialized.
* Controller actions are bound for the template.

So the **controller’s context** is:

| Property                                    | Description                          |
|---------------------------------------------|--------------------------------------|
| `this.model`                                | The upload data fetched in the route |
| `this.api`, `this.router`                   | Services injected via `@service`     |
| `this.deleteUpload()`, `this.shareUpload()` | Actions available to the template    |

---

### 3️⃣ Template renders

```
app/templates/uploads/show.hbs
```

The template is rendered using the controller as `this`:

```gjs
<template>
  <section>
    <h1>{{this.model.filename}}</h1>
    <Uploads::Details @upload={{this.model}} @onShare={{this.shareUpload}} />
  </section>
</template>
```

→ The controller’s `this.model` becomes the `@upload` argument for the component.

---

### 4️⃣ Component hierarchy takes over

```
app/components/uploads/details.gjs
```

Each component:

1. Receives arguments (`@upload`, `@onShare`).
2. Renders its internal template.
3. Handles local state (e.g., `isSharing` dialog flag).
4. Sends actions upward via the passed-in handlers (`this.args.onShare(allies)`).

If the user opens a modal and confirms “share,” it runs:

```
details component → calls @onShare(allies)
             │
             ▼
controller.shareUpload(allies)
             │
             ▼
api.shareUpload(model.id, allies)
```

---

### 5️⃣ API service handles network calls

```
app/services/api.js
```

* Centralized place for all `fetch` calls.
* Handles headers, credentials, JSON parsing, errors.

Returns a `Promise` → controller or component awaits it → tracked properties update → Ember re-renders automatically.

---

### 6️⃣ Data and rendering loop

```
Route (loads model)
  ↓
Controller (stores model, defines actions)
  ↓
Template (renders UI, passes data down)
  ↓
Component (renders subviews, calls actions up)
  ↓
Service (makes API requests)
```

After a fetch/POST succeeds:

* The controller may update `this.model` or transition.
* Ember automatically re-renders anything bound to tracked properties or `this.model`.

---

## 🧩 Visual Diagram (text version)

```
┌──────────────────────────────────────────────┐
│                Browser Request               │
│              /uploads/123                    │
└──────────────────────────────────────────────┘
                     │
                     ▼
┌──────────────────────────────────────────────┐
│ Ember Router (router.js)                     │
│   finds matching route: uploads.show         │
└──────────────────────────────────────────────┘
                     │
                     ▼
┌──────────────────────────────────────────────┐
│ Route: app/routes/uploads/show.js            │
│   beforeModel()                              │
│   model(params)  →  api.getUpload(id)        │
│   setupController(controller, model)         │
└──────────────────────────────────────────────┘
                     │
                     ▼
┌──────────────────────────────────────────────┐
│ Controller: app/controllers/uploads/show.js  │
│   this.model = upload                        │
│   actions: deleteUpload, shareUpload         │
└──────────────────────────────────────────────┘
                     │
                     ▼
┌──────────────────────────────────────────────┐
│ Template: app/templates/uploads/show.hbs     │
│   uses controller as context                 │
│   passes @upload + @onShare → component      │
└──────────────────────────────────────────────┘
                     │
                     ▼
┌──────────────────────────────────────────────┐
│ Component: app/components/uploads/details.gjs│
│   receives @upload, @onShare                 │
│   shows info; handles dialog state           │
│   calls this.args.onShare(allies)            │
└──────────────────────────────────────────────┘
                     │
                     ▼
┌──────────────────────────────────────────────┐
│ Service: app/services/api.js                 │
│   shareUpload(id, allies) → fetch POST       │
│   returns JSON result                        │
└──────────────────────────────────────────────┘
                     │
                     ▼
┌──────────────────────────────────────────────┐
│ Controller updates model / transitions       │
│ Ember re-renders templates automatically     │
└──────────────────────────────────────────────┘
```

---

### 🧭 TL;DR mental model

* **Router**: decides *which* routes are active.
* **Route**: decides *what data* to load.
* **Controller**: holds that data and defines *what can happen* (actions).
* **Template**: lays out *how* the data looks.
* **Component**: small, focused pieces of UI that send *actions up*.
* **Service**: reusable singletons that perform *cross-cutting tasks* (fetch, auth, storage).

---

## Appendix - Route Hierarchy

Route hierarchy diagram based on the router we’ve been talking about, i.e.:

```js
// app/router.js
this.route('uploads', function () {
    this.route('show', { path: '/:upload_id' });
    this.route('edit', { path: '/:upload_id/edit' });
});
```

Remember: Ember also creates `uploads.index` for you automatically. So the real route tree looks like this:

```text
application
└── uploads
    ├── uploads.index   (URL: /uploads)
    └── uploads.show    (URL: /uploads/:upload_id)
```

Now let’s line that up with the files mentioned earlier.

---

### 1. Route hierarchy → URLs → files

```text
application
└── uploads
    ├── uploads.index
    └── uploads.show
```

Becomes:

| Route name      | URL pattern           | Route file                          | Controller file                           | Template file                     |
| --------------- | --------------------- | ----------------------------------- | ----------------------------------------- |-----------------------------------|
| `application`   | `/` (root)            | `app/routes/application.js` (opt)   | `app/controllers/application.js` (opt)    | `app/templates/application.gjs`   |
| `uploads`       | `/uploads` (parent)   | `app/routes/uploads.js` (opt)       | `app/controllers/uploads.js` (opt)        | `app/templates/uploads.hbs`       |
| `uploads.index` | `/uploads`            | `app/routes/uploads/index.js` (opt) | `app/controllers/uploads/index.js` (rare) | `app/templates/uploads/index.gjs` |
| `uploads.show`  | `/uploads/:upload_id` | `app/routes/uploads/show.js`        | `app/controllers/uploads/show.js`         | `app/templates/uploads/show.gjs`  |

Notes:

* `uploads` is the parent; it renders `{{outlet}}`.
* `uploads.index` and `uploads.show` both render **into** that outlet.

---

### 2. Render flow (with outlets)

Think of the templates like nested boxes:

```text
application.gjs
└── uploads.gjs
    ├── (if URL is /uploads)        → render uploads/index.gjs here
    └── (if URL is /uploads/123)    → render uploads/show.gjs here
```

Or in pseudo-GJS:

```gjs
// app/templates/application.gjs
<template>
  <ApplicationShell>
    {{outlet}}
  </ApplicationShell>
</template>
```

```gjs
// app/templates/uploads.gjs
<template>
  <section class="uploads-layout">
    {{outlet}}  {{!-- child route goes here --}}
  </section>
</template>
```

If URL is `/uploads` → Ember puts `uploads/index.hbs` in that inner `{{outlet}}`.
If URL is `/uploads/123` → Ember puts `uploads/show.hbs` in that inner `{{outlet}}`.

---

### 3. Where the component fits

In the “show” case:

```text
application.hbs
└── uploads.hbs
    └── uploads/show.hbs
        └── <Uploads::Details ...>
```

So the full stack for `/uploads/123` is:

1. `app/router.js` matched `uploads.show`
2. `app/routes/uploads/show.js` loaded the model
3. `app/controllers/uploads/show.js` got the model + actions
4. `app/templates/uploads/show.hbs` rendered the page
5. `app/components/uploads/details.gjs` rendered the detailed view

---
