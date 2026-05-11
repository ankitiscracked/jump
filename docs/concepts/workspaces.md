# Workspaces

A workspace is the fundamental unit of work in jmp. Each workspace is a directory containing project files and a `.jmp/` metadata directory that tracks snapshots, manifests, and configuration.

## Local Configuration

Each workspace stores its config at `.jmp/config.json` as a `ProjectConfig`:

| Field               | Description                                      |
|---------------------|--------------------------------------------------|
| `project_id`        | ID of the project this workspace belongs to      |
| `workspace_id`      | Unique ID for this workspace                     |
| `workspace_name`    | Human-readable name                              |
| `base_snapshot_id`  | The snapshot this workspace was forked from       |
| `current_snapshot_id` | The most recent snapshot (auto-derived if empty)|
| `mode`              | Workspace mode string, usually empty or `local`  |
| `api_url`           | Legacy API URL field                             |

The `.jmp/` directory also contains:
- `snapshots/` -- snapshot `.meta.json` files
- `manifests/` -- manifest JSON files (keyed by content hash)
- `.gitignore` -- ignores snapshots/manifests from git

Implementation: `internal/config/config.go` (`ProjectConfig`, `Load`, `Save`, `InitAt`).

## Workspace Registry

All workspaces within a project are tracked in a project-level registry at `.jmp/workspaces/<workspace-id>.json` (stored in the project root's `.jmp/` directory). Each file contains a single workspace's metadata:

```json
{
  "workspace_id": "ws-abc123",
  "workspace_name": "feature-x",
  "path": "/path/to/feature-x",
  "base_snapshot_id": "snap-123",
  "current_snapshot_id": "snap-456",
  "created_at": "2025-01-01T00:00:00Z"
}
```

The registry enables cross-workspace commands like `jmp drift` and `jmp merge` to locate other workspaces by name. Per-workspace files avoid concurrent write conflicts when multiple workspaces operate in parallel.

Implementation: `internal/store/` (`Store`, `WorkspaceInfo`, `RegisterWorkspace`, `FindWorkspaceByName`).

## Lifecycle

### Init (`jmp workspace init`)

Creates a new workspace in the current directory:
1. Creates `.jmp/` with `config.json`, `snapshots/`, `manifests/`
2. Creates `.jmpignore` with default patterns if missing
3. Optionally creates an initial snapshot (`--no-snapshot` to skip)
4. Registers the workspace in the project-level workspace registry

Implementation: `cmd/jmp/commands/workspace.go` (`runInit`), `internal/config/config.go` (`InitAt`).

### Create (`jmp workspace create`)

Creates a new workspace directory under the current project folder. Forks from the source workspace's latest snapshot with all files copied. The new workspace is linked to the project's ID and placed as a subdirectory.

Implementation: `cmd/jmp/commands/workspace.go` (`newWorkspaceCreateCmd`).

## Main Workspace

Each project can designate one workspace as the "main" workspace. This is
stored in the project config as `main_workspace_id`. The main workspace serves
as the default comparison target for `jmp drift` when no workspace argument is
given.

Set via: `jmp workspace set-main [workspace-name]`

Implementation: `cmd/jmp/commands/workspace.go` (`runSetMain`).

## Workspace Status

`jmp status` displays current workspace info including name, ID, path, mode, snapshots, upstream, and change summary.

`jmp info workspaces` lists all workspaces for the current project from the local registry, showing name, path, role (main), and drift summary.

`jmp info workspace [name|id]` shows details for a specific workspace including snapshots, upstream, and project info.

Implementation: `cmd/jmp/commands/status.go` (`runStatus`), `cmd/jmp/commands/info.go` (`runInfoWorkspaces`, `runInfoWorkspace`).

## Storage

- Blobs: `.jmp/blobs/` (project-scoped, shared across workspaces under the same project)
- Global config: `~/.config/jmp/` (respects `XDG_CONFIG_HOME`) — agent preferences and author identity

Implementation: `internal/config/config.go` (`GetBlobsDir`, `GetBlobsDirAt`, `GetGlobalConfigDir`).

## Related Docs

- [Snapshots](snapshots.md) -- the immutable state records within workspaces
- [Drift](drift.md) -- detecting changes within a workspace
- [Merge](merge.md) -- merging changes between workspaces
- [Sync](sync.md) -- syncing workspace state through a backend
