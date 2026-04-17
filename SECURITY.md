# Security Policy

## Reporting a Vulnerability

Please **do not** open a public GitHub issue for security vulnerabilities.

Instead, email **sean@lidy.me** with:

- A description of the vulnerability and its potential impact
- Steps to reproduce or a minimal proof-of-concept
- Any suggested mitigations if you have them

You can expect an acknowledgement within 48 hours and a resolution timeline within 7 days depending on severity.

## Scope

golor is a pure computation library with no network access, file I/O, or external service calls. Security issues are most likely to involve incorrect output that could silently mislead accessibility checks (e.g. a contrast ratio reported as passing when it does not). These are treated as bugs with high priority.
