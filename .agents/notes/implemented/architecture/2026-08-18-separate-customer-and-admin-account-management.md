# Agent Note: Separate customer and administrator account management

Status: implemented

## Problem

Customer accounts and management accounts share authentication records but serve different operational purposes. Listing them together makes invitation metrics ambiguous and increases the risk of changing an administrator while performing customer support work.

## Decision

Accounts continue to use the same `users` table and authentication flow. Management lists are separated by effective management access: an administrator has the `super_admin` role or a role granting `referral:read`, `system:rbac`, or `system:*`; all other accounts are customers.

`/admin/users` shows customers and `/admin/administrators` shows administrators. Filtering happens in the backend query before counting and pagination. Role assignment remains the authority for moving an account between these views.

## Alternatives considered

**Separate administrator table:** This duplicates authentication, status, password, and audit behavior and complicates promotion of an existing customer.

**Frontend-only filtering:** This produces incorrect totals and incomplete pages because filtering happens after server pagination.

**Any assigned role means administrator:** Some roles may be non-management roles, so classification follows effective management permissions instead.

## Consequences

Operations have clear customer and administrator workspaces without duplicating identities. Changing relevant role grants can move an account between lists, so role changes and user changes must remain audited.

## Verification

- Repository tests cover server-side customer and administrator filtering.
- Frontend typecheck and production builds cover both management routes.
