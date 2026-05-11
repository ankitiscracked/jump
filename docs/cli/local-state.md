# Local State Layout

How `jmp` stores state on disk. Source: `internal/config/config.go`, `internal/store/`.

## Per-Workspace: `.jmp/`

Created by `jmp workspace init` or `jmp project init` in the workspace root.

```
.jmp/
  config.json       # Workspace configuration
  lock              # Exclusive flock for workspace-level locking
  author.json       # Optional project-level author identity (name, email)
  snapshots/        # Snapshot metadata files
    <id>.meta.json  # SnapshotMeta: ID, ManifestHash, Parents, Author, CreatedAt, Message
  manifests/        # Manifest files (file tree hashes)
    <hash>.json     # Manifest: map of relative paths to content hashes
  blobs/            # Content-addressed blob storage (project-scoped)
    <sha256-hash>   # Raw file content keyed by SHA-256 hash
  workspaces/       # Project-level workspace registry (project root only)
    <ws-id>.json    # Per-workspace metadata (name, path, snapshot IDs)
  tasks/            # Project-level task/session metadata
    <task-id>.json  # Task: snapshots and files grouped by unit of work
  stat-cache.json   # Stat cache for accelerating manifest generation
  current-task.json # Workspace-local pointer to the active task, if any
  export/           # Git export state (created by `jmp git export`)
    git-map.json    # Snapshot ID <-> git commit SHA mapping
  events/           # Append-only project-local coordination events
    <id>-<type>.json
  .gitignore        # Auto-generated, ignores .jmp internals
```

### `config.json` fields

```json
{
  "project_id": "uuid",
  "workspace_id": "uuid",
  "workspace_name": "my-workspace",
  "base_snapshot_id": "uuid",
  "current_snapshot_id": "uuid",
  "api_url": "https://...",
  "mode": ""
}
```

Defined as `ProjectConfig` in `internal/config/config.go`. The `fork_snapshot_id` field is deprecated and migrated to `base_snapshot_id` on load.

## Global Config: `~/.config/jmp/`

Respects `XDG_CONFIG_HOME`. Contains user-level state shared across all projects.

```
~/.config/jmp/
  agents.json       # Preferred agent configuration
  author.json       # Global author identity (name, email)
```

### `agents.json` structure

```json
{
  "preferred_agent": "claude"
}
```

Source: `internal/agent/agent.go`.

## `.jmpignore`

Located at the workspace root (alongside `.jmp/`). Created automatically by `jmp workspace init` if not present. See [ignore.md](ignore.md) for pattern syntax. A default set of patterns is embedded into the binary from `internal/ignore/default.jmpignore`.

## Project root detection

`config.FindProjectRoot()` walks up from the current directory looking for a `.jmp/` directory. This determines which workspace context commands operate in.
