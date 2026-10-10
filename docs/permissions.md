# Permissions (RBAC)

`k8sec` makes three Kubernetes API calls, all on `secrets` in the core API group:

| Action in k8sec | API verb |
| --- | --- |
| Show the list | `list` |
| Save an edited Secret | `update` |
| Add a new Secret (`a`) | `create` |

It never deletes Secrets and doesn't need `get`.

## Minimal namespaced Role

This grants exactly what `k8sec` needs in one namespace (`demo` here):

```yaml
apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  name: k8sec-editor
  namespace: demo
rules:
  - apiGroups: [""]
    resources: ["secrets"]
    verbs: ["list", "update", "create"]
---
apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
  name: k8sec-editor
  namespace: demo
subjects:
  # Bind to a user, group or ServiceAccount. Example: a ServiceAccount.
  - kind: ServiceAccount
    name: k8sec-user
    namespace: demo
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: Role
  name: k8sec-editor
```

Read-only browsing? Drop `update` and `create`: you'll be able to list and open Secrets,
but saving will fail.

This Role was tested with `k8sec v0.0.2` using a ServiceAccount token kubeconfig: listing,
editing and creating in `demo` all worked, and running without `-n demo` failed with
`cannot list resource "secrets" ... in the namespace "default"`, as expected.

## Check your access

```bash
kubectl auth can-i list   secrets -n demo
kubectl auth can-i update secrets -n demo
kubectl auth can-i create secrets -n demo
```

## A kubeconfig for a ServiceAccount

`k8sec` only reads `~/.kube/config`, so to use a ServiceAccount, put its context there and
make it current:

```bash
TOKEN=$(kubectl -n demo create token k8sec-user)
kubectl config set-credentials k8sec-user --token="$TOKEN"
kubectl config set-context k8sec --cluster="$(kubectl config view --minify -o jsonpath='{.clusters[0].name}')" --user=k8sec-user
kubectl config use-context k8sec
k8sec -n demo
```

Tokens from `kubectl create token` expire (one hour by default).
