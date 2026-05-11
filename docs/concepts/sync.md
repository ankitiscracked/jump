# Sync

`jmp sync` reconciles local snapshot state with a configured backend.

Backends are configured at project level:

```bash
jmp backend set git
jmp backend set github owner/repo
```

The current backend implementations use Git export/import as the transport.
When local and backend heads diverge, sync uses the same three-way merge
machinery as `jmp merge`.

## Commands

- `jmp backend status`: show backend config.
- `jmp backend push`: push local snapshot export to backend.
- `jmp pull`: pull latest backend state.
- `jmp sync`: pull, merge if needed, then push.

`jmp sync` conflict flags:

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

Implementation: `cmd/jmp/commands/sync.go`, `cmd/jmp/commands/backend.go`,
`cmd/jmp/commands/cloud_merge.go`, `internal/backend/`.
