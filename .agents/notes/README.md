# Agent Notes

Agent Notes are lightweight decision records written for humans and coding agents. They explain why a non-trivial product or engineering choice exists, which alternatives were rejected, and how the result is verified.

## Layout

Use `{lifecycle}/{category}/YYYY-MM-DD-topic.md`:

- Lifecycle: `proposed` or `implemented`.
- Category: `architecture`, `feature`, `bug-fix`, `process`, or `testing`.

Do not create a note for routine formatting, isolated copy changes, or other decisions already obvious from code. Create or update one when a change affects module ownership, referral/Credit semantics, authentication or authorization boundaries, persistent data, public API contracts, deployment routing, or team-wide workflow.

## Proposed format

```markdown
# Agent Note: Title

Status: proposed

## Problem
## Proposal
## Alternatives considered
## Acceptance criteria
## Risks
```

## Implemented format

```markdown
# Agent Note: Title

Status: implemented

## Problem
## Decision
## Alternatives considered
## Consequences
## Verification
```

Write observable facts and rationale, not chat transcripts or hidden reasoning. Link relevant code and other notes with relative paths. Keep implemented notes current when paths, names, or delivered mechanics change; use a new cross-linked note when reversing the decision itself.
