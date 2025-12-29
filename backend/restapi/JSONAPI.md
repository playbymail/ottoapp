# JSON:API view files

Put JSON:API “view” structs like `ClanDocumentView` in the **server (transport) layer**, not in `domains`.

Why:

* **Domain stays transport-agnostic.** Your `domains` package should model *business concepts and rules* (Document, permissions, IDs, invariants), not how one particular API serializes them. JSON:API tags (`jsonapi:"attr,..."`) are pure transport concerns.
* **Views change more often than domain.** Field names like `"document-name"` vs `"documentName"`, resource type strings, `iso8601` formatting, link shapes, relationship inclusion, etc. evolve with clients and spec choices.
* **You’ll likely have multiple representations.** Same domain object might be rendered differently for `/user.documents` vs `/admin.documents`, or for list vs detail. Keeping those as transport DTOs avoids contaminating domain with presentation variants.

A clean split I’ve used in Go looks like:

* `domains/`
    * `type ClanDocument struct { … }`
    * `type DocumentType string`
    * permission logic / policy helpers (or methods on a “permission set”)
* `server/http/jsonapi/` (or `internal/transport/jsonapi/`)
    * `type ClanDocumentView struct { … jsonapi tags … }`
    * `func ToClanDocumentView(d domains.ClanDocument, perms domains.DocPerms) ClanDocumentView`

One nuance: if you have *non-JSON:API* “view models” that are still useful across multiple transports (say, HTTP + CLI + a future gRPC), those still usually don’t belong in `domains`; they belong in an `app`/`service` layer (sometimes called `usecase`), as DTOs without JSON tags. Then the JSON:API view is one more mapping on top.

So: **server package** for the JSON:API-tagged struct; **domains** for the underlying business types and permission rules.
