# Backend k6 checks

The script exercises authenticated user and administrator read paths without
creating users, referrals, or Credit ledger entries. Referral-code creation is
safe to repeat because the endpoint returns the account's existing code.

Run a short local smoke test:

```bash
ADMIN_EMAIL='admin@example.com' \
ADMIN_PASSWORD='replace-me' \
k6 run tests/k6/backend-smoke.js
```

Use a separate ordinary-user account when available:

```bash
ADMIN_EMAIL='admin@example.com' \
ADMIN_PASSWORD='replace-me' \
USER_EMAIL='user@example.com' \
USER_PASSWORD='replace-me' \
VUS=5 DURATION=1m \
k6 run tests/k6/backend-smoke.js
```

Configuration:

- `BASE_URL`: backend origin, default `http://localhost:8999`
- `ADMIN_EMAIL` / `ADMIN_PASSWORD`: required administrator credentials
- `USER_EMAIL` / `USER_PASSWORD`: optional; defaults to the administrator
- `VUS`: user-scenario virtual users, default `2`
- `DURATION`: scenario duration, default `20s`

The test fails when checks drop below 99%, failed requests reach 1%, p95 exceeds
750 ms, or p99 exceeds 1.5 s.
