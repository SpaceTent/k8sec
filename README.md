# k8sec

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENCE)
[![Go version](https://img.shields.io/github/go-mod/go-version/SpaceTent/k8sec)](go.mod)
[![Latest tag](https://img.shields.io/github/v/tag/SpaceTent/k8sec)](https://github.com/SpaceTent/k8sec/tags)

**A simple terminal UI for editing Kubernetes Secrets.**

`k8sec` lists the Secrets in a namespace, lets you pick one with a searchable list, and
opens its decoded values in your editor as plain `KEY=VALUE` lines. Save and quit, and the
Secret is updated in the cluster. No more `kubectl get secret -o yaml`, base64-decoding,
editing, re-encoding and applying by hand.

![k8sec listing the Secrets in the demo namespace](docs/images/k8sec-list.png)

## Contents

- [Features](#features)
- [Requirements](#requirements)
- [Installation](#installation)
- [Quick start](#quick-start)
- [Usage](#usage)
- [Keybindings](#keybindings)
- [How it works](#how-it-works)
- [Limitations](#limitations)
- [Security notes](#security-notes)
- [Project layout](#project-layout)
- [Contributing](#contributing)
- [License](#license)

## Features

- **Interactive list** of every Secret in a namespace, built with
  [Bubble Tea](https://github.com/charmbracelet/bubbletea).
- **Fuzzy filtering**: press `/` and type part of a name.
- **Edit in your own editor** (`$EDITOR`, default `nano`) as `.env`-style `KEY=VALUE` lines,
  with values already base64-decoded.
- **Change detection**: if you quit without changing anything, nothing is sent to the
  cluster and `k8sec` prints `no changes`.
- **Create new Secrets** (type `Opaque`) from the same screen with `a`.

| Filter | Edit |
| --- | --- |
| ![Filtering the list by typing "post"](docs/images/k8sec-filter.png) | ![The postgres-credentials Secret open in nano as KEY=VALUE lines](docs/images/k8sec-edit.png) |

## Requirements

- Access to a Kubernetes cluster, configured in **`~/.kube/config`**. `k8sec` always reads
  that file and uses its **current context**. (`$KUBECONFIG` and `--context` are not
  supported yet; switch with `kubectl config use-context` first.)
- Permission to `list`, `update` and `create` Secrets in the namespace you use. See
  [docs/permissions.md](docs/permissions.md) for a minimal Role.
- A terminal text editor. Set `EDITOR` (for example `vim`, `nvim`, `nano`); if unset,
  `k8sec` runs `nano`.
- To build: **Go 1.26+** and `git` (and `make` if you want to use the Makefile).

## Installation

There are no pre-built binaries yet, so build from source:

```bash
git clone https://github.com/SpaceTent/k8sec.git
cd k8sec
make build          # produces ./bin/k8sec with the version and build time embedded
```

Without `make`:

```bash
go build -o bin/k8sec .
```

Then put the binary on your `PATH`, for example:

```bash
sudo install -m 0755 bin/k8sec /usr/local/bin/k8sec
```

> `go install github.com/SpaceTent/k8sec@latest` does **not** work at the moment, because
> `go.mod` declares the module as `k8sec` rather than `github.com/SpaceTent/k8sec`.

Check the build:

```console
$ ./bin/k8sec -h
k8sec: TUI-based Kubernetes Secret editor
Version: v0.0.2
BuildTime: 20261010-160808

Usage of k8sec:
  -n string
    	namespace to use (overrides NAMESPACE env var, defaults to 'default')
```

## Quick start

```bash
# Make sure kubectl points at the right cluster
kubectl config current-context

# Browse and edit Secrets in the "demo" namespace
k8sec -n demo
```

1. Use the arrow keys (or `j`/`k`) to highlight a Secret and press **Enter**.
2. Your editor opens with one `KEY=VALUE` line per key:

   ```text
   POSTGRES_USER=app
   POSTGRES_DB=app
   POSTGRES_PASSWORD=change-me
   ```

3. Change a value, save and quit. `k8sec` updates the Secret:

   ```text
   Secret postgres-credentials updated.
   ```

   If you quit without changing anything, it prints `no changes` and leaves the cluster
   alone.

A step-by-step walkthrough, including creating a Secret, is in
[docs/usage.md](docs/usage.md).

## Usage

```text
k8sec [-n namespace]
```

| Option | Description |
| --- | --- |
| `-n <namespace>` | Namespace to work in. Takes priority over `NAMESPACE`. |
| `NAMESPACE` env var | Namespace to use when `-n` isn't given. |
| `EDITOR` env var | Editor command to run. Defaults to `nano`. |
| `-h` | Show help, version and build time. |

If neither `-n` nor `NAMESPACE` is set, `k8sec` uses `default`. It does **not** read the
namespace set on your kubeconfig context.

```bash
k8sec                       # namespace "default"
k8sec -n my-namespace       # explicit namespace (recommended)
NAMESPACE=my-namespace k8sec
EDITOR=vim k8sec -n my-namespace
```

## Keybindings

In the list:

| Key | Action |
| --- | --- |
| `↑`/`↓` or `k`/`j` | Move the selection |
| `Enter` | Edit the highlighted Secret |
| `/` | Start filtering by name |
| `a` | Add a new Secret |
| `?` | Show more key help (paging, jump to start/end) |
| `q` or `Esc` | Quit without editing |
| `Ctrl+C` | Quit |

While filtering:

| Key | Action |
| --- | --- |
| Type | Narrow the list (fuzzy match on the Secret name) |
| `Enter` | Edit the highlighted (top) match straight away |
| `Esc` | Cancel the filter |

> **Heads-up:** `a` is caught before the filter sees it, so typing a name that contains
> the letter `a` while filtering switches to *add new Secret* mode. Filter on other letters
> for now (for example `post` rather than `postgres-cred`).

## How it works

1. **Connect**: builds a client from `~/.kube/config` (current context) using `client-go`.
2. **List**: fetches all Secrets in the namespace and shows them in a Bubble Tea list.
3. **Edit**: writes the selected Secret's data to a temporary file (`$TMPDIR/ksec-*.env`,
   mode `0600`) as `KEY=VALUE` lines, then runs `$EDITOR <file>`.
4. **Compare**: after the editor exits, re-reads the file and compares it with the
   original data.
5. **Update**: if anything changed, replaces the Secret's `data` and calls `Update`.
   Otherwise prints `no changes`. The temp file is deleted either way.
6. **Create** (`a`): asks for a name and a namespace (defaulting to the current one),
   opens an empty file in your editor, and creates an `Opaque` Secret from the lines you
   enter. An empty file creates nothing.

There's a diagram and more detail in [docs/architecture.md](docs/architecture.md).

## Limitations

Please read these before using `k8sec` on important Secrets.

> [!WARNING]
> **Multi-line values are not supported and can be lost.** The editor format is one
> `KEY=VALUE` per line. A value containing newlines (TLS certificates, private keys, JSON
> or YAML blobs) is split across lines; only the first line is kept on the way back, and
> because that differs from the original, `k8sec` **writes the truncated value to the
> cluster even if you didn't edit anything**. Use `kubectl` for those Secrets.

- **Order:** keys appear in random order each time (they come from a Go map).
- **Comments:** there's no comment syntax. A line like `# note=x` becomes a key called
  `# note`. Lines without `=` and blank lines are ignored.
- **Whitespace:** keys and values are taken exactly as typed, including spaces around `=`.
- **Removing a key:** delete its line; the key is removed on update.
- **Binary values** that aren't valid text will not round-trip cleanly.
- **Kubeconfig:** only `~/.kube/config` and its current context are used.
- **Editor command:** `EDITOR` must be a single program name or path; arguments (such as
  `code --wait`) aren't supported. Use a wrapper script if you need them.
- **Secret names:** the add prompt mentions `_`, but Kubernetes names must be lowercase
  letters, digits, `-` and `.`; the API rejects anything else.
- **Errors:** connection or permission problems currently exit with a Go panic and stack
  trace. The first line has the real message, for example
  `secrets is forbidden: User "..." cannot list resource "secrets" ... in the namespace "default"`.
- **Secret type:** new Secrets are always `Opaque`. Editing keeps the existing type.

## Security notes

- Kubernetes Secrets are **base64-encoded, not encrypted**. `k8sec` shows the decoded
  values, so anyone who can see your screen can read them.
- While you edit, the values are in a temporary file readable only by you (`0600`) in
  your system temp directory. It's removed when `k8sec` exits normally, but your editor
  may keep its own swap/backup files (for example Vim's `.swp`). Configure your editor
  accordingly.
- `k8sec` uses whatever permissions your kubeconfig grants. Prefer a namespaced Role with
  only `list`, `update` and `create` on `secrets` ([docs/permissions.md](docs/permissions.md)).

To report a vulnerability, see [SECURITY.md](SECURITY.md).

## Project layout

```text
.
├── main.go       # Flags, kube client, list/edit/create flow, editor and .env handling
├── model.go      # Bubble Tea model: key handling for the secret list
├── item.go       # List item (title = Secret name, description = namespace)
├── Makefile      # build / fmt / release helpers
├── docs/         # Usage guide, permissions, architecture, screenshots
└── LICENCE
```

## Contributing

Issues and pull requests are welcome. Please read [CONTRIBUTING.md](CONTRIBUTING.md) and
the [Code of Conduct](CODE_OF_CONDUCT.md) first.

## License

`k8sec` is licensed under the MIT Licence. See the [LICENCE](LICENCE) file for details.
