# jmp

jmp is a local-first CLI for parallel coding-agent work. It gives agents and
humans durable workspace snapshots, cheap workspace forks, drift detection,
three-way merge, rollback, and Git export.

## Install From Source

```bash
go build -o ./bin/jmp ./cmd/jmp
```

## Quick Start

The CLI help starts with the same happy path:

```bash
jmp --help
```

Create a project with a `main` workspace:

```bash
jmp project create myproject
cd myproject/main
```

Capture work:

```bash
jmp snapshot -m "Initial version"
```

Create a parallel workspace:

```bash
jmp workspace create agent-auth
cd ../agent-auth
```

Compare and merge:

```bash
jmp drift main
jmp merge main
```

Export to Git:

```bash
jmp git export
```

## Daily Commands

- `jmp task start` starts a named unit of work for an agent/human.
- `jmp snapshot` captures current workspace state.
- `jmp task status` shows snapshots and files for the active task.
- `jmp workspace create` forks a peer workspace.
- `jmp status` shows local changes since the latest snapshot.
- `jmp drift` compares this workspace against another workspace.
- `jmp merge` performs a three-way merge.
- `jmp task finish` closes the active unit of work.
- `jmp restore` restores files from a snapshot.
- `jmp events` shows project-local coordination events.
- `jmp git export` turns snapshot history into Git commits.

## Docs

- [CLI commands](docs/cli/commands.md)
- [Local state](docs/cli/local-state.md)
- [Ignore patterns](docs/cli/ignore.md)
- [Agents](docs/cli/agents.md)
- [Packaged agent skill](.agents/skills/jmp-agent-workflow/SKILL.md)
- [Snapshots](docs/concepts/snapshots.md)
- [Workspaces](docs/concepts/workspaces.md)
- [Drift](docs/concepts/drift.md)
- [Merge](docs/concepts/merge.md)
- [Git integration](docs/concepts/git-integration.md)

## Tests

```bash
go test ./...
./tests/e2e/run.sh
```
