# Agent Note: Separate user and admin applications

Status: implemented

## Problem

The Referral product serves two audiences with different navigation, authorization, and deployment needs. Treating the customer experience and operations console as one application obscures ownership and makes root routing, admin base paths, and permission boundaries easy to break.

## Decision

The frontend is a pnpm workspace. `frontend/apps/web` owns `/` and `/ref/:code`; `frontend/apps/admin` owns `/admin/`. Shared HTTP contracts and browser client behavior live in `frontend/packages/api`, while reusable presentation primitives live in `frontend/packages/ui`. Deployment routes `/api/` to the Go backend and preserves independent SPA fallbacks for the two applications.

Authentication and authorization remain explicit per audience. A normal authenticated user does not gain console access without the required administrative role or permission.

## Alternatives considered

**One React application with an admin route:** This reduces initial configuration, but couples bundles, navigation, route guards, and release risk for two different audiences.

**Independent repositories:** This maximizes isolation but adds versioning and coordination overhead for shared API contracts that is not justified at the current product size.

## Consequences

Each application can evolve and deploy with a clear base path, while shared packages prevent duplicate API plumbing. Changes to shared packages must be checked against both applications, and nginx plus React Router base configuration must stay aligned.

## Verification

- `cd frontend && pnpm typecheck && pnpm build` checks both applications.
- `/` serves the user app and `/ref/:code` supports invitation registration.
- `/admin/` serves the console with assets resolved beneath the admin base path.
- `/api/` reaches the backend without being captured by an SPA fallback.
