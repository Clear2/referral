---
name: referral-release-check
description: Check whether Referral backend, user web, admin console, database migrations, configuration, and deployment routing are ready to release. Use before handoff, deployment, merge, or when asked to verify the whole product.
---

# Referral Release Check

Produce an evidence-based release verdict. Do not deploy, commit, or mutate external systems unless explicitly requested.

## Checks

1. Inspect the current diff and classify changes by backend, web, admin, shared packages, migration, configuration, and deployment.
2. Confirm no secrets, certificates, production config, raw tokens, or generated Ent edits are newly tracked.
3. Run backend formatting/linting when available and `go test ./...`. Treat race-enabled `make test` as the stronger final check when its dependencies are available.
4. From `frontend/`, run `pnpm typecheck` and `pnpm build` so both apps are checked.
5. When routing or authentication changed, manually verify:
   - `/` serves the user application.
   - `/ref/:code` supports invited registration.
   - `/admin/` serves the operations console and its assets under the correct base.
   - `/api/` reaches the backend.
   - expired/invalid sessions reach the appropriate login page without loops.
   - logout exists and clears the session.
6. When schema changes, confirm a forward migration exists, generated Ent code matches the schema, and example configuration/docs are updated without real credentials.
7. When referral or Credit behavior changes, exercise success, duplicate/retry, unauthorized, and disabled-user paths; verify displayed aggregates against the ledger/source data.

## Verdict

Return one of `READY`, `READY WITH CAVEATS`, or `NOT READY`. List executed commands and results, blockers ordered by severity, manual checks, and explicitly unverified items. Do not claim readiness based only on compilation.
