# CLI Command Reference

Reference for the `fst` CLI commands in `cmd/fst/commands`.

## Projects

| Command | Description |
|---------|-------------|
| `fst project init [name]` | Wrap or initialize a project folder |
| `fst project create <name>` | Create a new local project with a `main` workspace |

Flags:

- `project init`: `--project-id`, `--keep-name`, `--force`
- `project create`: `--no-snapshot`, `--force`, `--path`

## Workspaces

| Command | Description |
|---------|-------------|
| `fst workspace init [project-name]` | Initialize a workspace inside an existing project folder |
| `fst workspace create <workspace-name>` | Fork a new workspace from current/main/source workspace |
| `fst workspace set-main [workspace]` | Set the project main workspace |

Flags:

- `workspace init`: `--workspace, -w`, `--no-snapshot`, `--force`
- `workspace create`: `--from`, `--backend` (`auto`, `clone`, `copy`)

## Snapshots

| Command | Description |
|---------|-------------|
| `fst snapshot` / `fst snap` | Capture current workspace state |
| `fst log` | Show snapshot history |
| `fst restore [files...]` | Restore files from a snapshot |

Flags:

- `snapshot`: `--message, -m`, `--agent-message`
- `log`: `--limit, -n`, `--all, -a`, `--graph, -g`
- `restore`: `--to`, `--to-base`, `--dry-run`

## Drift, Diff, Merge

| Command | Description |
|---------|-------------|
| `fst status` | Show workspace status and local changes |
| `fst drift [workspace]` | Compare workspaces using DAG merge-base |
| `fst diff [workspace] [file...]` | Show line-level differences |
| `fst merge <workspace>` | Three-way merge from another workspace |
| `fst dag` | Show project-wide snapshot DAG |
| `fst events` | Show project-local coordination events |
| `fst watch` | Stream project-local coordination events |

Flags:

- `status`: `--json`
- `drift`: `--json`, `--agent-summary`, `--no-dirty`
- `diff`: `--context, -C`, `--no-color`, `--names-only`
- `merge`: `--manual`, `--theirs`, `--ours`, `--dry-run`, `--agent-summary`, `--no-pre-snapshot`, `--force`, `--abort`
- `dag`: `--limit, -n`
- `events`: `--since`, `--json`
- `watch`: `--since`, `--interval`

Current event types: `workspace_created`, `snapshot_created`, `merge_started`,
`merge_completed`, `restore_completed`.

Merge conflict modes:

- Default: use preferred coding agent if available, with manual fallback.
- `--manual`: write conflict markers.
- `--theirs`: accept source version for conflicts.
- `--ours`: keep current version for conflicts.

## Info

| Command | Description |
|---------|-------------|
| `fst info` | Show context-aware workspace/project info |
| `fst info workspaces` / `fst info ws` | List project workspaces |
| `fst info workspace [name|id]` | Show one workspace |
| `fst info project` | Show current project |

Flags:

- `info`: `--json`
- `info workspace`: `--json`
- `info project`: `--json`

## Git And Backends

| Command | Description |
|---------|-------------|
| `fst git export` | Export snapshot DAG to Git commits |
| `fst git import <repo-path>` | Import from a Git repo exported by `fst` |
| `fst github export <owner>/<repo>` | Export to a GitHub repository |
| `fst github import <owner>/<repo>` | Import from a GitHub repository |
| `fst backend set git` | Use local Git export as backend |
| `fst backend set github <owner/repo>` | Use GitHub remote as backend |
| `fst backend push` | Push backend state |
| `fst pull` | Pull latest backend state |
| `fst sync` | Sync local and backend state |

Flags:

- `git export`: `--branch, -b`, `--include-dirty`, `--message, -m`, `--init`, `--rebuild`
- `git import`: `--branch, -b`, `--workspace, -w`, `--project, -p`, `--rebuild`
- `github export`: `--branch, -b`, `--include-dirty`, `--message, -m`, `--init`, `--rebuild`, `--remote`, `--create`, `--private`, `--push-all`, `--force-remote`, `--no-gh`
- `github import`: `--branch, -b`, `--workspace, -w`, `--project, -p`, `--rebuild`, `--no-gh`
- `backend set github`: `--create`, `--private`, `--remote`, `--force-remote`
- `sync`: `--manual`, `--theirs`, `--ours`

## Agents

| Command | Description |
|---------|-------------|
| `fst agents` | List detected coding agents |
| `fst agents list` / `fst agents ls` | List known agents and availability |
| `fst agents set-preferred [name]` | Set preferred coding agent |

## Configuration

| Command | Description |
|---------|-------------|
| `fst config` | Interactive author setup |
| `fst config --global` | Interactive global author setup |
| `fst config set name "Jane Doe"` | Set project author name |
| `fst config set email "jane@example.com"` | Set project author email |
| `fst config get [key]` | Show resolved author config |

Flags:

- `config`: `--global`
- `config set`: `--global`

## Advanced

| Command | Description |
|---------|-------------|
| `fst edit <snapshot>` | Edit snapshot message |
| `fst drop <snapshot>` | Drop snapshot from current chain |
| `fst squash <from>..<to>` | Squash a linear range |
| `fst rebase <from>..<to> --onto <snapshot>` | Rebase a linear range |
| `fst gc` | Prune unreachable snapshots/blobs |
| `fst ui` | Open terminal UI |
| `fst version` | Print build/version info |

Deprecated:

- `fst conflicts <workspace-path>`: use `fst merge --dry-run`.
