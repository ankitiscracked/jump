# JMP Agent Workflow

Use this skill when working inside a jmp (`jmp`) workspace, when the user asks
you to use `jmp`, or when several agents/humans may edit the same project.

## Goal

Let the coding agent do the coding. Let `jmp` provide durable checkpoints,
parallel workspaces, drift detection, merge, and Git export.

## Start Of Work

1. Check context:

```bash
jmp status --json
jmp task status --json
```

If `jmp task status --json` says there is no active task and you will edit
files, start one:

```bash
jmp task start "<short task name>"
```

If you are in a project root and need an isolated workspace:

```bash
jmp workspace create <agent-or-task-name>
cd <agent-or-task-name>
jmp task start "<short task name>"
```

## During Work

Use normal coding tools for code search, edits, tests, and reasoning.

After each coherent milestone, create a checkpoint:

```bash
jmp snapshot -m "<what changed>"
```

Before merging or handing off, check drift:

```bash
jmp drift main --json
jmp diff main --names-only
```

If `main` is not the right target, use the workspace name given by
`jmp info workspaces`.

## Coordination

To inspect nearby activity:

```bash
jmp events --json
```

To watch while running a long task in another terminal:

```bash
jmp watch
```

Treat new events as awareness, not orchestration. If another workspace changes
the same files you are editing, run `jmp drift <workspace>` and decide whether
to merge, re-plan, or ask the user.

## Finish

Run tests or checks appropriate to the repository, then snapshot the final
state:

```bash
jmp snapshot -m "<final checkpoint>"
jmp task finish --summary "<what is done>"
```

For human review in Git/GitHub:

```bash
jmp git export --init
```

Exported commits include `Jmp-Snapshot`, workspace, task, and agent metadata
trailers when available.

## Hook Ideas

If the host agent supports hooks, keep them advisory and cheap:

- On session start: run `jmp status --json` and `jmp task status --json`.
- After file edits: optionally run `jmp status --json`.
- Before a long test/build: run `jmp snapshot -m "before <command>"` if there
  are meaningful changes.
- On session end: run `jmp snapshot -m "<summary>"` and `jmp task finish` when
  the task is complete.

Do not wrap the agent process with `jmp`. The agent should call `jmp` directly
at the points where snapshots, drift, tasks, or export matter.
