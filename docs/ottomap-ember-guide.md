# OttoMap Ember 6.8 Data Flow and Route Structure Guide

This guide explains how data moves through an Ember application (v6.8), using OttoMap’s `/uploads/show/:upload_id` page as an example.

---

## 1. Overview

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

## 2. File Layout Example

```
app/
├── routes/
│   └── uploads/
│       └── show.js
├── controllers/
│   └── uploads/
│       └── show.js
├── templates/
│   └── uploads/
│       └── show.hbs
├── components/
│   └── uploads/
│       └── details.gjs
├── services/
│   └── api.js
└── tests/
    └── acceptance/
        └── uploads-show-test.js
```

---

## 3. Flow of Control

### 1️⃣ Router

`app/router.js` maps URLs to routes:

```js
Router.map(function () {
  this.route('uploads', function () {
    this.route('show', { path: '/:upload_id' });
  });
});
```

Visiting `/uploads/123` activates the route chain:

```
application → uploads → uploads.show
```

---

### 2️⃣ Route

`app/routes/uploads/show.js`

```js
import Route from '@ember/routing/route';
import { service } from '@ember/service';

export default class UploadsShowRoute extends Route {
  @service api;

  async model(params) {
    return this.api.getUpload(params.upload_id);
  }

  setupController(controller, model) {
    super.setupController(controller, model);
    controller.upload = model;
  }
}
```

Responsibilities:
- Fetch data via the API service.
- Set up the controller’s model.

---

### 3️⃣ Controller

`app/controllers/uploads/show.js`

```js
import Controller from '@ember/controller';
import { service } from '@ember/service';
import { action } from '@ember/object';

export default class UploadsShowController extends Controller {
  @service api;
  @service router;

  @action
  async deleteUpload() {
    if (!confirm(`Delete ${this.model.filename}?`)) return;
    await this.api.deleteUpload(this.model.id);
    this.router.transitionTo('uploads.index');
  }

  @action
  async shareUpload(allies) {
    await this.api.shareUpload(this.model.id, allies);
    alert('Upload shared successfully.');
  }
}
```

Responsibilities:
- Hold the route’s model.
- Define actions to mutate or delete data.
- Coordinate with services for side effects.

---

### 4️⃣ Template

`app/templates/uploads/show.hbs`

```hbs
<section class="p-4 space-y-6">
  <header class="flex justify-between items-center">
    <h1 class="text-xl font-semibold">Upload: {{this.model.filename}}</h1>
    <Ui::Button {{on "click" this.deleteUpload}}>Delete</Ui::Button>
  </header>

  <Uploads::Details @upload={{this.model}} @onShare={{this.shareUpload}} />
</section>
```

Responsibilities:
- Define layout and pass arguments to components.
- Bind actions from the controller to UI elements.

---

### 5️⃣ Component

`app/components/uploads/details.gjs`

```gjs
<template>
  <div class="rounded-lg border p-4 bg-white">
    <Uploads::FileInfo @upload={{@upload}} />
    <Ui::Button {{on "click" this.openShareDialog}}>Share</Ui::Button>

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

  @action openShareDialog() { this.isSharing = true; }
  @action cancelShare() { this.isSharing = false; }

  @action async confirmShare(allies) {
    await this.args.onShare(allies); // calls controller action
    this.isSharing = false;
  }
}
```

Responsibilities:
- Handle presentation and local state.
- Call actions passed from parents (“actions up”).

---

### 6️⃣ Service

`app/services/api.js`

```js
import Service from '@ember/service';

export default class ApiService extends Service {
  baseUrl = '/api';

  async getUpload(id) {
    return this.#request(`/uploads/${id}`);
  }

  async deleteUpload(id) {
    return this.#request(`/uploads/${id}`, { method: 'DELETE' });
  }

  async shareUpload(id, allies) {
    return this.#request(`/uploads/${id}/share`, {
      method: 'POST',
      body: JSON.stringify({ allies }),
    });
  }

  async #request(path, options = {}) {
    const res = await fetch(`${this.baseUrl}${path}`, {
      headers: { 'Content-Type': 'application/json' },
      credentials: 'include',
      ...options,
    });
    if (!res.ok) throw new Error(`${res.status} ${res.statusText}`);
    return res.status === 204 ? null : res.json();
  }
}
```

Responsibilities:
- Perform all backend communication.
- Normalize responses and handle errors.

---

## 4. Sequence Diagram

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

## 5. Key Rules for OttoMap

| Layer | Purpose | Fetch? | Typical File |
|--------|----------|--------|--------------|
| **Route** | Load data from backend | ✅ Yes (via `api`) | `app/routes/...` |
| **Controller** | Manage model and actions | ✅ Maybe | `app/controllers/...` |
| **Template** | Bind data and layout | ❌ No | `app/templates/...` |
| **Component** | Present UI, bubble actions up | ❌ No (except arg-driven) | `app/components/...` |
| **Service** | Handle network requests, state | ✅ Yes | `app/services/api.js` |

---

## 6. TL;DR

- Routes load data.  
- Controllers hold data and actions.  
- Templates render the UI.  
- Components display reusable UI and send actions up.  
- Services perform network requests.

---

**Data flows down; actions go up.**
