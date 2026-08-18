Run the `referral-release-check` project skill against the current change.

Inspect the changed scope, run backend tests plus frontend typecheck/build, check secrets and migrations, and manually verify relevant user/admin/auth/referral routes. Return `READY`, `READY WITH CAVEATS`, or `NOT READY` with evidence and explicit unverified items. Do not deploy or commit unless requested.
