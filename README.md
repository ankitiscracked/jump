# Jump (`fst`)

Jump is a local-first CLI for parallel coding-agent work. It gives agents and
humans durable workspace snapshots, cheap workspace forks, drift detection,
three-way merge, rollback, and Git export.

The binary is still named `fst`.

## Install From Source

```bash
go build -o ./bin/fst ./cmd/fst
```

## Quick Start

Create a project with a `main` workspace:

```bash
fst project create myproject
cd myproject/main
```

Capture work:

```bash
fst snapshot -m "Initial version"
```

Create a parallel workspace:

```bash
fst workspace create agent-auth
cd ../agent-auth
```

Compare and merge:

```bash
fst drift main
fst merge main
```

Export to Git:

```bash
fst git export
```

## Daily Commands

- `fst snapshot` captures current workspace state.
- `fst workspace create` forks a peer workspace.
- `fst status` shows local changes since the latest snapshot.
- `fst drift` compares this workspace against another workspace.
- `fst merge` performs a three-way merge.
- `fst restore` restores files from a snapshot.
- `fst events` shows project-local coordination events.
- `fst git export` turns snapshot history into Git commits.

## Docs

- [CLI commands](docs/cli/commands.md)
- [Local state](docs/cli/local-state.md)
- [Ignore patterns](docs/cli/ignore.md)
- [Agents](docs/cli/agents.md)
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
