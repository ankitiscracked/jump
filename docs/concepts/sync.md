# Sync

`fst sync` reconciles local snapshot state with a configured backend.

Backends are configured at project level:

```bash
fst backend set git
fst backend set github owner/repo
```

The current backend implementations use Git export/import as the transport.
When local and backend heads diverge, sync uses the same three-way merge
machinery as `fst merge`.

## Commands

- `fst backend status`: show backend config.
- `fst backend push`: push local snapshot export to backend.
- `fst pull`: pull latest backend state.
- `fst sync`: pull, merge if needed, then push.

`fst sync` conflict flags:

- `--manual`: write conflict markers.
- `--theirs`: take backend version for conflicts.
- `--ours`: keep local version for conflicts.

## Divergence

When histories diverge:

1. Find merge base.
2. Load local and backend manifests.
3. Apply non-conflicting backend changes.
4. Resolve conflicts by selected mode.
5. Create a merge snapshot with both parents.
6. Push merged state through the backend.

Implementation: `cmd/fst/commands/sync.go`, `cmd/fst/commands/backend.go`,
`cmd/fst/commands/cloud_merge.go`, `internal/backend/`.
