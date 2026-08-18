# Agent Note: Connect RBAC menus and API resources

Status: implemented

## Problem

Roles, permissions, menus, and API resources were persisted, but menu assignments did not drive the admin navigation or route access. Permission-to-API and permission-to-menu relationships could be read from the RBAC snapshot but not edited in the admin UI, and API mappings did not participate in authorization.

## Decision

The current-user access response now includes enabled menu records. Effective menus are the union of menus assigned directly to enabled roles and menus assigned through those roles' enabled permissions. The admin root loader checks the requested route against that effective menu set and the shared admin navigation renders the same records. `super_admin` remains a route-access fallback.

Permission editors can assign API resources and menus. Existing hard-coded permission checks remain authoritative; an enabled API resource mapping can additionally satisfy a route check when the current user owns one of its enabled permissions. This preserves existing safe defaults when no resource mapping exists.

A migration creates the four built-in admin menu records and maps the referral dashboard to `referral:read` and account/RBAC pages to `system:rbac`.

## Alternatives considered

- Replace all hard-coded permission checks with database mappings: rejected because an incomplete or accidental configuration could expose administrative endpoints.
- Use role-menu relationships only: rejected because permission bundles should be able to carry their required navigation resources.
- Keep hard-coded frontend navigation: rejected because persisted menu assignments would continue to be informational rather than enforceable.

## Consequences

Menu configuration now affects both navigation visibility and client route entry. Backend authorization remains the security boundary. Custom frontend routes still require a corresponding compiled React Router route; adding a database menu does not dynamically load arbitrary components.

## Verification

- `go test -count=1 ./internal/modules/rbac`
- `go test ./internal/modules/...`
- `cd frontend && pnpm typecheck && pnpm build`
