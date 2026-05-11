# FST Agent Workflow

Use this skill when working inside a Jump (`fst`) workspace, when the user asks
you to use `fst`, or when several agents/humans may edit the same project.

## Goal

Let the coding agent do the coding. Let `fst` provide durable checkpoints,
parallel workspaces, drift detection, merge, and Git export.

## Start Of Work

1. Check context:

```bash
fst status --json
fst task status --json
```

If `fst task status --json` says there is no active task and you will edit
files, start one:

```bash
fst task start "<short task name>"
```

If you are in a project root and need an isolated workspace:

```bash
fst workspace create <agent-or-task-name>
cd <agent-or-task-name>
fst task start "<short task name>"
```

## During Work

Use normal coding tools for code search, edits, tests, and reasoning.

After each coherent milestone, create a checkpoint:

```bash
fst snapshot -m "<what changed>"
```

Before merging or handing off, check drift:

```bash
fst drift main --json
fst diff main --names-only
```

If `main` is not the right target, use the workspace name given by
`fst info workspaces`.

## Coordination

To inspect nearby activity:

```bash
fst events --json
```

To watch while running a long task in another terminal:

```bash
fst watch
```

Treat new events as awareness, not orchestration. If another workspace changes
the same files you are editing, run `fst drift <workspace>` and decide whether
to merge, re-plan, or ask the user.

## Finish

Run tests or checks appropriate to the repository, then snapshot the final
state:

```bash
fst snapshot -m "<final checkpoint>"
fst task finish --summary "<what is done>"
```

For human review in Git/GitHub:

```bash
fst git export --init
```

Exported commits include `Fst-Snapshot`, workspace, task, and agent metadata
trailers when available.

## Hook Ideas

If the host agent supports hooks, keep them advisory and cheap:

- On session start: run `fst status --json` and `fst task status --json`.
- After file edits: optionally run `fst status --json`.
- Before a long test/build: run `fst snapshot -m "before <command>"` if there
  are meaningful changes.
- On session end: run `fst snapshot -m "<summary>"` and `fst task finish` when
  the task is complete.

Do not wrap the agent process with `fst`. The agent should call `fst` directly
at the points where snapshots, drift, tasks, or export matter.
