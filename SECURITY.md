# Security Policy

`k8sec` reads and writes Kubernetes Secrets, so we take security reports seriously.

## Supported versions

There are no binary releases yet. Fixes land on `master` and are tagged as the next
`v0.0.x`. Please check the issue still exists on the latest `master` before reporting.

## Reporting a vulnerability

**Please do not open a public issue, discussion or pull request for security problems.**

Report privately using GitHub's
[private vulnerability reporting](https://github.com/SpaceTent/k8sec/security/advisories/new)
(**Security** tab → **Report a vulnerability**). If that isn't available, email
**mark@ackroyd.net** with `[k8sec security]` in the subject.

Please include:

- What the issue is and what an attacker could do with it.
- Steps to reproduce (use dummy Secrets, never real values).
- `k8sec` version or commit, OS, and Kubernetes version.

We'll acknowledge your report as soon as we can, keep you informed while we fix it, and
credit you in the release notes unless you prefer to stay anonymous.

## In scope

- Secret values leaking outside the temporary edit file (for example world-readable
  files, logs or crash output).
- The temporary file not being removed, or being created with weak permissions.
- `k8sec` writing to a Secret, namespace or cluster other than the one selected.

## Out of scope

- Values visible on your own screen or in your editor's swap/backup files. Configure your
  editor (for example `:set noswapfile` in Vim) if this matters to you.
- Access granted by your kubeconfig/RBAC. `k8sec` can do whatever your credentials allow.
