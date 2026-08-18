---
name: referral-feature
description: Implement or change Referral product features across the Go backend, shared API package, user web app, and admin app. Use for referral links, invitation registration, Credit rewards, authentication, users, RBAC, analytics, Ent schema, migrations, or cross-application behavior.
---

# Referral Feature

Deliver an end-to-end product change without breaking the separation between the public user experience and the operations console.

## Workflow

1. Read the root `AGENTS.md`, the relevant module, route, shared package, and tests before editing.
2. Trace the behavior from HTTP contract through Controller → Service → Repository and into the owning UI. Confirm whether the change belongs to `web`, `admin`, or both.
3. Define invariants before implementation. For reward changes, state who earns Credit, on which event, whether it is idempotent, and how duplicate or failed registration behaves.
4. Keep backend request types in their owning module. Change `ent/schema` and add a migration when persistence changes; regenerate Ent and never hand-edit generated files.
5. Put reusable HTTP types and client logic in `frontend/packages/api`; put truly shared presentation primitives in `frontend/packages/ui`. Keep product flows inside their owning app.
6. Preserve routing boundaries: `/` and `/ref/:code` are user-facing; `/admin/` is the console; `/api/` is backend traffic. A 401/403 caused by an expired user session must lead the owning app to its login flow without creating a redirect loop.
7. Add focused tests at the changed seam, then run the proportional checks described below.
8. For a non-trivial architectural or behavioral decision, add or update an Agent Note under `.agents/notes`.

## Product rules to verify

- Referral registration must bind the invited user to the inviter exactly once.
- Credit awards must be transactional or otherwise idempotent and auditable.
- User and admin authorization are separate concerns; an authenticated user is not automatically an administrator.
- Disabling a user must block protected operations consistently.
- UI totals such as successful invitations and accumulated Credit must come from backend data, not placeholders.
- Do not expose passwords, tokens, referral internals, or privileged fields in logs or responses.

## Verification

Run targeted tests while developing. Before handoff, normally run:

```bash
(cd backend && go test ./internal/modules/...)
(cd frontend && pnpm typecheck && pnpm build)
```

Also verify the affected happy path and at least one authorization, duplicate, or failure path. If a command cannot run because infrastructure is unavailable, report the exact unverified behavior.
