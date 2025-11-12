# Ember Data Flow Pattern (DDAU)

**Service → Route → Controller → Template → Component**

| Layer | Responsibility | Fetch? |
|--------|----------------|--------|
| **Service** | Centralizes backend communication (`fetch`, headers, error handling) | ✅ Yes |
| **Route** | Calls the API service, loads data in `model()`, and provides it to the controller | ✅ Yes |
| **Controller** | Holds the `@model`, defines page-level actions (`save`, `delete`, etc.), and transitions between routes | ⚙️ No direct fetch |
| **Template** | Binds `this.model` and controller actions to components as arguments | ❌ |
| **Component** | Renders UI, reacts to user input, and calls passed-in actions upward (`@onSave`, `@onCancel`) | ❌ |

**Rule of thumb:**  
> Routes call the API service → Controllers wire `@model` into templates → Components render it and `@actions` bubble it back up.

**Only the API service and Ember Simple Auth should ever call `fetch()`.**
