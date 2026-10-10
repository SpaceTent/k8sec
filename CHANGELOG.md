# Changelog

All notable changes to this project are documented here.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/), and this
project uses [Semantic Versioning](https://semver.org/) tags (pre-1.0). Entries for
existing tags were reconstructed from the git history.

## [Unreleased]

### Added
- CONTRIBUTING, CODE_OF_CONDUCT, SECURITY, CHANGELOG, issue and PR templates.
- Expanded README with screenshots, `docs/usage.md`, `docs/permissions.md` and
  `docs/architecture.md`.

## [v0.0.2] - 2026-04-11

### Removed
- Compiled binary removed from the repository; `bin/` is now git-ignored.

## [v0.0.1] - 2026-04-11

### Added
- Initial release: list Secrets in a namespace in a Bubble Tea TUI, fuzzy filter, edit as
  `KEY=VALUE` in `$EDITOR`, change detection, create new `Opaque` Secrets, `-n` flag and
  `NAMESPACE` env var.

[Unreleased]: https://github.com/SpaceTent/k8sec/compare/v0.0.2...HEAD
[v0.0.2]: https://github.com/SpaceTent/k8sec/compare/v0.0.1...v0.0.2
[v0.0.1]: https://github.com/SpaceTent/k8sec/releases/tag/v0.0.1
