---
name: referral-code-review
description: Review Referral repository changes for correctness, regressions, security, data integrity, and release risk across the Go backend, shared API package, user web app, admin console, migrations, RBAC, authentication, referrals, and Credit. Use when asked to review a diff, branch, commit, pull request, working tree, implementation, or pre-merge change without modifying it.
---

# Referral Code Review

Review changed behavior, not formatting preferences. Report only defects that are concrete, reachable, and worth fixing before merge.

## Establish the review scope

1. Read the root `AGENTS.md` and relevant Agent Notes.
2. Use the comparison point supplied by the user. Otherwise inspect the working tree and determine a sensible merge base with the repository's main branch.
3. List changed files and read the complete affected functions, callers, tests, migrations, and contracts. Do not review isolated diff lines without surrounding behavior.
4. If version-control context is unavailable, state the exact files or implementation surface reviewed and the resulting limitation.

Keep the review read-only. Do not edit, commit, migrate a database, or change external state unless the user separately asks for fixes.

## Trace affected behavior

For each changed flow, trace the first entry point through its real dependencies:

- HTTP: nginx/path → Gin middleware → Controller → Service → Repository → PostgreSQL/Redis/message system.
- Frontend: React Router route/loader → shared API client → backend contract → loading, success, unauthorized, forbidden, empty, and failure states.
- Persistence: Ent schema → generated model expectations → SQL migration/up-down behavior → existing data compatibility.

Confirm that displayed totals and permissions come from real backend data rather than placeholders or frontend-only state.

## Review invariants

### Authentication and authorization

- Keep user and administrator authorization separate.
- Distinguish 401 from 403 and disabled users; prevent redirect loops.
- Require backend authorization even when menus or controls are hidden in the UI.
- Verify RBAC role, permission, menu, and API-resource associations affect the intended enforcement boundary.
- Prevent self-demotion, privilege escalation, insecure bootstrap behavior, and accidental `super_admin` bypasses.

### Referral and Credit

- Bind an invitee to an inviter at most once.
- Prevent self-referral and duplicate or replayed registration rewards.
- Keep referral creation, Credit balance changes, and ledger writes transactional or otherwise idempotent.
- Ensure failed registration cannot leave a referral or reward behind.
- Reconcile dashboard totals with persisted rows and pagination/filter semantics.

### Data and migrations

- Never accept edits to generated Ent files as the source of truth.
- Require forward and rollback migrations for schema changes; protect existing rows and constraints.
- Check foreign-key deletion behavior, uniqueness, nullable transitions, seed idempotency, and migration hashes.
- Reject destructive rollback logic that can delete user-created data matching a seed heuristic.

### Frontend and deployment

- Preserve `/`, `/ref/:code`, `/admin/`, and `/api/` ownership and production basename behavior.
- Check shared API envelope types, cookie credentials, session refresh behavior, and safe `next` redirects.
- Verify Chinese/English copy in rendered UI and programmatic channels such as share text, email bodies, titles, and accessibility labels.
- Check mobile form behavior, loading states, empty data, accessible controls, and route-level SEO/noindex boundaries when touched.
- Ensure static assets, prerender paths, nginx fallbacks, and environment-derived canonical URLs work in the deployed layout.

## Validate findings

Before reporting a defect:

1. Identify the exact changed line and the first incorrect state or contract.
2. Prove a realistic input, actor, route, or failure path reaches it.
3. Check whether validation, middleware, a database constraint, or a caller already prevents it.
4. Run the narrowest relevant read-only test when practical. Use broader checks only in proportion to risk.
5. Do not report legacy defects unrelated to the reviewed change unless the change makes them newly reachable.

Typical checks include:

```bash
(cd backend && go test ./internal/modules/...)
(cd frontend && pnpm typecheck && pnpm build)
```

Do not treat passing builds as proof of behavioral correctness.

## Report format

Lead with findings ordered by severity:

- `P0`: immediate security breach, irreversible loss, or system-wide outage.
- `P1`: likely production failure, authorization bypass, or material data corruption.
- `P2`: user-visible regression or incorrect behavior on a realistic path.
- `P3`: smaller but actionable correctness, maintainability, or operability defect.

For every finding include:

- concise title with severity;
- clickable file and line reference;
- triggering scenario and observable impact;
- why existing safeguards do not prevent it;
- smallest repair direction, without implementing it.

After findings, list assumptions or unverified runtime conditions and the checks run. If no defects qualify, say so explicitly and name remaining test gaps. Avoid praise, change summaries, style nits, and speculative warnings that lack a reachable failure path.
