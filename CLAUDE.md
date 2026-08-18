# Referral AI Guide

Read and follow [AGENTS.md](AGENTS.md) as the repository authority.

Project skills are shared from `.agents/skills` and exposed to Claude-compatible tools through `.claude/skills`:

- `referral-feature` for end-to-end product changes.
- `referral-investigate` for evidence-driven defect diagnosis.
- `referral-release-check` before handoff or release.

For non-trivial decisions, read `.agents/notes/README.md` and record the rationale in the appropriate lifecycle and category. Never place credentials, tokens, production configuration, or private user data in prompts, notes, logs, or generated files.
