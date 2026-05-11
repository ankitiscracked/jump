# CLI Command Reference

Reference for the `jmp` CLI commands in `cmd/jmp/commands`.

## Projects

| Command | Description |
|---------|-------------|
| `jmp project init [name]` | Wrap or initialize a project folder |
| `jmp project create <name>` | Create a new local project with a `main` workspace |

Flags:

- `project init`: `--project-id`, `--keep-name`, `--force`
- `project create`: `--no-snapshot`, `--force`, `--path`

## Workspaces

| Command | Description |
|---------|-------------|
| `jmp workspace init [project-name]` | Initialize a workspace inside an existing project folder |
| `jmp workspace create <workspace-name>` | Fork a new workspace from current/main/source workspace |
| `jmp workspace set-main [workspace]` | Set the project main workspace |

Flags:

- `workspace init`: `--workspace, -w`, `--no-snapshot`, `--force`
- `workspace create`: `--from`, `--backend` (`auto`, `clone`, `copy`)

## Snapshots

| Command | Description |
|---------|-------------|
| `jmp snapshot` / `jmp snap` | Capture current workspace state |
| `jmp log` | Show snapshot history |
| `jmp restore [files...]` | Restore files from a snapshot |

Flags:

- `snapshot`: `--message, -m`, `--agent-message`
- `log`: `--limit, -n`, `--all, -a`, `--graph, -g`
- `restore`: `--to`, `--to-base`, `--dry-run`

## Tasks

Tasks group multiple snapshots into one agent/human unit of work. They are a
thin CLI primitive over snapshots; agents can keep using `jmp snapshot` while
the active task records which snapshots and files belong together.

| Command | Description |
|---------|-------------|
| `jmp task start <name>` | Start a task in the current workspace |
| `jmp task status [task-id]` | Show the current or named task |
| `jmp task finish` | Finish the current task |
| `jmp task list` | List project tasks |

Flags:

- `task status`: `--json`
- `task finish`: `--summary`
- `task list`: `--json`

## Drift, Diff, Merge

| Command | Description |
|---------|-------------|
| `jmp status` | Show workspace status and local changes |
| `jmp drift [workspace]` | Compare workspaces using DAG merge-base |
| `jmp diff [workspace] [file...]` | Show line-level differences |
| `jmp merge <workspace>` | Three-way merge from another workspace |
| `jmp dag` | Show project-wide snapshot DAG |
| `jmp events` | Show project-local coordination events |
| `jmp watch` | Stream project-local coordination events |

Flags:

- `status`: `--json`
- `drift`: `--json`, `--agent-summary`, `--no-dirty`
- `diff`: `--context, -C`, `--no-color`, `--names-only`
- `merge`: `--manual`, `--theirs`, `--ours`, `--dry-run`, `--agent-summary`, `--no-pre-snapshot`, `--force`, `--abort`
- `dag`: `--limit, -n`
- `events`: `--since`, `--json`
- `watch`: `--since`, `--interval`

Current event types: `workspace_created`, `snapshot_created`, `task_started`,
`task_finished`, `merge_started`, `merge_completed`, `restore_completed`.

Merge conflict modes:

- Default: use preferred coding agent if available, with manual fallback.
- `--manual`: write conflict markers.
- `--theirs`: accept source version for conflicts.
- `--ours`: keep current version for conflicts.

## Info

| Command | Description |
|---------|-------------|
| `jmp info` | Show context-aware workspace/project info |
| `jmp info workspaces` / `jmp info ws` | List project workspaces |
| `jmp info workspace [name|id]` | Show one workspace |
| `jmp info project` | Show current project |

Flags:

- `info`: `--json`
- `info workspace`: `--json`
- `info project`: `--json`

## Git And Backends

| Command | Description |
|---------|-------------|
| `jmp git export` | Export snapshot DAG to Git commits |
| `jmp git import <repo-path>` | Import from a Git repo exported by `jmp` |
| `jmp github export <owner>/<repo>` | Export to a GitHub repository |
| `jmp github import <owner>/<repo>` | Import from a GitHub repository |
| `jmp backend set git` | Use local Git export as backend |
| `jmp backend set github <owner/repo>` | Use GitHub remote as backend |
| `jmp backend push` | Push backend state |
| `jmp pull` | Pull latest backend state |
| `jmp sync` | Sync local and backend state |

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
| `jmp agents` | List detected coding agents |
| `jmp agents list` / `jmp agents ls` | List known agents and availability |
| `jmp agents set-preferred [name]` | Set preferred coding agent |

## Configuration

| Command | Description |
|---------|-------------|
| `jmp config` | Interactive author setup |
| `jmp config --global` | Interactive global author setup |
| `jmp config set name "Jane Doe"` | Set project author name |
| `jmp config set email "jane@example.com"` | Set project author email |
| `jmp config get [key]` | Show resolved author config |

Flags:

- `config`: `--global`
- `config set`: `--global`

## Advanced

| Command | Description |
|---------|-------------|
| `jmp edit <snapshot>` | Edit snapshot message |
| `jmp drop <snapshot>` | Drop snapshot from current chain |
| `jmp squash <from>..<to>` | Squash a linear range |
| `jmp rebase <from>..<to> --onto <snapshot>` | Rebase a linear range |
| `jmp gc` | Prune unreachable snapshots/blobs |
| `jmp ui` | Open terminal UI |
| `jmp version` | Print build/version info |

Deprecated:

- `jmp conflicts <workspace-path>`: use `jmp merge --dry-run`.
