---
name: referral-investigate
description: Systematically diagnose Referral system defects involving login, cookies, 401/403 responses, referral registration, Credit accounting, RBAC, user state, React Router, API calls, or production path routing. Use when behavior is broken or its root cause is unknown.
---

# Referral Investigate

Find the first broken boundary and establish the root cause before changing code.

## Investigation loop

1. Restate the expected behavior, actual behavior, affected URL, actor type, and environment.
2. Reproduce with the smallest path. Capture the HTTP status, response body, relevant server log, and browser console/network evidence without exposing secrets.
3. Trace one request end to end: nginx/path routing → React Router/client → shared API client → Gin middleware/controller → service → repository/database.
4. Test competing hypotheses with read-only checks or focused tests. Do not treat the last stack frame as the cause without evidence.
5. Identify the first incorrect state or contract, explain why it produces the symptom, and check adjacent flows for the same defect.
6. If the user requested a fix, make the smallest coherent repair and add a regression test. If they requested diagnosis only, stop after the evidence-backed finding.

## High-value checks

- For 401/403: distinguish unauthenticated, expired/invalid session, disabled user, missing role, and missing permission. Confirm cookie domain, path, `Secure`, SameSite, and credentials mode.
- For referrals: inspect referral-code resolution, inviter/invitee uniqueness, transaction boundaries, retry behavior, and registration rollback.
- For Credit: reconcile ledger rows against displayed totals; look for duplicate event processing and non-atomic updates.
- For `/admin/`: confirm Vite/React Router base paths, asset URLs, nginx fallback, and that `/` still serves the user app.
- For frontend stack traces: find the first application frame and inspect generated route configuration before editing dependencies.

## Verification

Run the narrowest regression test first, then the affected application checks. Report the root cause, evidence, changed files, and any remaining unverified conditions.
