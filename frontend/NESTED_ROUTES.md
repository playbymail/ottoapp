# Nested Routes imply Nested UI

In Ember, the structure of your `router.js` determines not just the URLs of your application, but also how your User Interface (UI) is nested.

## The Concept

If you have a route nested inside another in `router.js`:

```javascript
Router.map(function() {
  this.route('posts', function() {
    this.route('show', { path: '/:post_id' });
  });
});
```

Ember assumes the UI should also be nested. When you visit `/posts/1`:
1. Ember renders the **Parent** template (`posts.gjs`).
2. It looks for an `{{outlet}}` inside `posts.gjs`.
3. It renders the **Child** template (`posts/show.gjs`) into that outlet.

## The Common Pitfall

A common mistake is to put the "List of Posts" directly in the **Parent** template (`posts.gjs`).

If you do this:
1. **Visiting `/posts`**: You see the list.
2. **Visiting `/posts/1`**: You see the list **AND** the single post detail below it (if you included an `{{outlet}}`), or you see just the list (if you forgot the `{{outlet}}`).

Usually, this is not what you want. You want the "Show" page to **replace** the "List" page.

## The Solution: The "Index" Route Pattern

To make sibling pages mutually exclusive (either show the List OR show the Detail, but never both), you use the **Index Route**.

1. **Parent (`posts.gjs`)**: Becomes a hollow shell. It only contains `{{outlet}}`. It stays active for all child routes.
2.  **Index (`posts/index.gjs`)**: Contains the List. This is rendered *only* when you are exactly at `/posts`.
3.  **Detail (`posts/show.gjs`)**: Contains the Detail. This is rendered *only* when you are at `/posts/1`.

## Example Directory Structure

Here is how you should structure your files to achieve this "List vs Detail" behavior:

```text
app/
├── router.js                 // Defines parent 'posts' and child 'show'
└── templates/
    └── posts.gjs             // PARENT: Contains only <template>{{outlet}}</template>
    └── posts/
        ├── index.gjs         // LIST: The list of posts. Renders at /posts
        └── show.gjs          // DETAIL: The single post. Renders at /posts/1
```

### In Summary

*   **Parent Template**: Shared layout (like a sidebar or just an empty slot).
*   **Index Template**: The "default" content for the parent URL.
*   **Sibling Template**: Content that replaces the Index content.
