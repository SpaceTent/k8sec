# Contributing to k8sec

Thanks for helping improve `k8sec`! Bug reports, docs, tests and features are all welcome.
By taking part you agree to follow the [Code of Conduct](CODE_OF_CONDUCT.md).

## Ground rules

- **Never use real Secrets** in issues, PRs, screenshots or tests. Use a throwaway cluster
  and dummy values.
- Security problems go through [SECURITY.md](SECURITY.md), not public issues.
- For anything bigger than a small fix, open an issue first so we can agree on the approach.

## Development setup

You need Go **1.26+** (see `go.mod`), `git`, and optionally `make`, `goimports` and
`govulncheck`.

```bash
git clone https://github.com/SpaceTent/k8sec.git
cd k8sec
go build ./...
go vet ./...
make build          # or: go build -o bin/k8sec .
./bin/k8sec -h
```

Optional tools used by the Makefile:

```bash
go install golang.org/x/tools/cmd/goimports@latest
go install golang.org/x/vuln/cmd/govulncheck@latest
```

## Trying your changes against a cluster

Use a disposable local cluster, never production. For example with
[kind](https://kind.sigs.k8s.io/):

```bash
kind create cluster --name k8sec-dev          # writes the context into ~/.kube/config
kubectl create namespace demo
kubectl -n demo create secret generic api-keys \
  --from-literal=STRIPE_KEY=sk_test_123 --from-literal=SENDGRID_KEY=SG.example
kubectl -n demo create secret generic redis-auth --from-literal=REDIS_PASSWORD=hunter2

./bin/k8sec -n demo

kind delete cluster --name k8sec-dev          # when you're done
```

minikube, k3d or Docker Desktop's Kubernetes work just as well. Remember that `k8sec`
always uses the **current context** in `~/.kube/config`, so double-check
`kubectl config current-context` before testing.

To test with limited permissions, use the Role in [docs/permissions.md](docs/permissions.md).

## Tests

There are no automated tests yet, and adding them is a great first contribution. Good
candidates are the pure functions in `main.go`: `secretToEnv`, `readEnv`, `hasChanges` and
`getNamespace`. Table-driven tests with the standard `testing` package are fine.

```bash
go test ./...
```

## Making a change

1. Fork the repo and branch from `master`: `git checkout -b fix/short-description`
   (or `feat/...`, `docs/...`).
2. Keep the change focused: one fix or feature per pull request.
3. Format and check:

   ```bash
   gofmt -s -w .
   goimports -w .     # if installed
   go vet ./...
   go test ./...
   govulncheck ./...  # if installed
   ```

4. Update `README.md` / `docs/` if behaviour, flags or keybindings change. If the TUI looks
   different, refresh the screenshots in `docs/images/` from a real run.
5. Add a line under **Unreleased** in [CHANGELOG.md](CHANGELOG.md).

### Make targets

| Target | What it does |
| --- | --- |
| `make build` | Builds `bin/k8sec` with `VERSION` (from `git describe --tags`) and `BUILDTIME` embedded |
| `make fmt` | `go mod tidy`, `goimports`, `gofmt -s` |
| `make commit` / `make release` | Maintainer shortcuts that commit, tag and **push**. Contributors should use normal git and pull requests instead. |

## Commit messages

Short, imperative subjects. A [Conventional Commits](https://www.conventionalcommits.org/)
prefix is preferred:

```text
fix: keep 'a' as text while filtering
feat: honour $KUBECONFIG
docs: document RBAC requirements
test: cover readEnv edge cases
```

## Pull requests

- Target `master` and fill in the PR template.
- Describe how you tested it (cluster type, editor, terminal).
- There's no CI yet, so please run the checks above locally.
