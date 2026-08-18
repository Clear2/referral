# Agent Note: Public SEO and prerendering

Status: implemented

## Problem

The user application was built as a single SPA fallback document. Route metadata was only applied after JavaScript loaded, so crawlers and social preview bots could not reliably read titles, descriptions, canonical URLs, or Open Graph data. Account, invitation, and admin routes also lacked an explicit indexing policy.

## Decision

The production site origin defaults to `https://referral.vivl.cc` and can be overridden with `VITE_SITE_URL`. The web application prerenders `/`, `/login`, and `/register`; the primary route is indexable while account routes emit `noindex`. Dynamic invitation URLs and `/admin/` are excluded through route metadata and `robots.txt`, and the admin document always emits `noindex`.

Shared SEO helpers own canonical, Open Graph, and Twitter descriptors. The web root owns application metadata, favicon/manifest links, and Organization/WebSite/SoftwareApplication JSON-LD. Static `robots.txt`, `sitemap.xml`, manifest icons, and a 1200×630 Open Graph image ship with the frontend public assets.

## Alternatives considered

- Client-only metadata: rejected because many preview bots do not execute the application.
- Full runtime SSR: deferred because the applications currently deploy as static SPAs and the public routes can be covered by build-time prerendering.
- Index invitation-code URLs: rejected to avoid exposing or duplicating user-specific invitation URLs in search results.

## Consequences

New public, indexable routes must be added to `prerender` and the sitemap. If the production origin changes, deployment should set `VITE_SITE_URL` and update the static sitemap/robots sitemap URL.

## Verification

- `cd frontend && pnpm typecheck`
- `cd frontend && pnpm build`
- Inspect prerendered `/`, `/login`, and `/register` HTML for canonical and robots metadata.
