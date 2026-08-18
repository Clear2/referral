# Agent Note: Shared frontend localization

Status: implemented

## Problem

The customer web application and admin console contained independently embedded Chinese copy. They had no common locale preference or language control, so adding English independently would allow the two deployed applications to drift.

## Decision

Own locale detection, persistence, English translations, and the language control in `frontend/packages/i18n`. Both application roots install the shared provider. The preference uses the `referral.locale` local-storage key, defaults to the browser language, and updates the document language attribute. Because the existing applications predate structured translation keys, the initial compatibility layer localizes rendered text and accessible attributes, including content mounted after API requests.

New product copy should be added to the shared catalog. Locale behavior does not belong in `@referral/api`, and product-specific layouts remain in their owning applications.

## Alternatives considered

- Separate dictionaries in each application: rejected because common authentication, navigation, and account terms would drift.
- A backend locale setting: deferred because this release only localizes browser UI and does not require account-level synchronization or API contract changes.
- Rewriting every route around translation keys in one change: deferred to keep this first migration reviewable; the shared compatibility layer gives existing routes immediate coverage.

## Consequences

The customer and admin applications share one saved language choice. Chinese remains the server-rendered fallback, while non-Chinese browsers default to English after hydration. Returning from English to Chinese reloads the current route to restore its original source copy.

## Verification

- `cd frontend && pnpm typecheck`
- `cd frontend && pnpm build`
