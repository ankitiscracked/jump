# Coding Agent Integration

How `jmp` detects and invokes coding agents for summaries and merge conflict resolution. Source: `internal/agent/agent.go`, `cmd/jmp/commands/agents.go`.

## Supported agents

| Name | Command | Invocation style | Description |
|------|---------|-----------------|-------------|
| `claude` | `claude` | `claude -p "<prompt>"` | Claude Code (Anthropic) |
| `codex` | `codex` | `codex exec "<prompt>"` | OpenAI Codex CLI |
| `amp` | `amp` | `amp -x "<prompt>"` | Amp CLI |
| `agent` | `agent` | `agent -p "<prompt>"` | Cursor Agent CLI |
| `gemini` | `gemini` | `gemini -p "<prompt>"` | Gemini CLI |
| `droid` | `droid` | `droid exec "<prompt>"` | Factory Droid CLI |

Defined as `KnownAgents` in `internal/agent/agent.go`.

## Agent detection

Agents are detected by checking if their command is on `PATH` using `exec.LookPath`. The `jmp agents list` command shows all known agents and whether each is currently available.

## Preferred agent

Set with `jmp agents set-preferred <name>`. Stored in `~/.config/jmp/agents.json`:

```json
{
  "preferred_agent": "claude"
}
```

When an agent is needed (for summaries or merges), the selection order is:

1. `--agent` flag (on commands that support it)
2. `JMP_AGENT` environment variable
3. Preferred agent from `agents.json`
4. First available agent from the known agents list

## Agent capabilities

Agents are invoked for four tasks:

### Snapshot summaries (`--agent-summary`)

`InvokeSummary(agentName, diffText)` -- generates a concise summary of changes in a snapshot. Used by `jmp snapshot --agent-summary`.

### Conflict summaries

`InvokeConflictSummary(agentName, conflictText)` -- summarizes detected conflicts between workspaces. Used by `jmp drift --agent-summary`.

### Drift summaries

`InvokeDriftSummary(agentName, driftText)` -- summarizes drift between workspaces. Used by `jmp drift --agent-summary`.

### Merge conflict resolution

`InvokeMerge(agentName, baseContent, oursContent, theirsContent, filePath)` -- resolves a three-way merge conflict. The agent receives base, ours, and theirs versions and returns a `MergeResult` containing:

- `Strategy` -- list of strings describing what the agent did
- `MergedCode` -- the resolved file content

The agent's output is parsed by looking for a `---MERGED CODE---` separator. Everything before it is treated as strategy explanation; everything after is the merged file content.

## Commands

```
jmp agents                      # List detected agents (runs list by default)
jmp agents list                 # Show all known agents with availability (alias: ls)
jmp agents set-preferred claude # Set preferred agent
```

## Commands that use agents

| Command | Flag | Purpose |
|---------|------|---------|
| `jmp snapshot` | `--agent-summary`, `--agent` | Generate snapshot summary |
| `jmp drift` | `--agent-summary` | Summarize drift/conflicts |
| `jmp merge` | `--agent-summary` | Summarize merge; agent mode resolves conflicts |
| `jmp pull` | `--agent-summary` | Summarize pulled changes |
| `jmp sync` | `--agent-summary` | Summarize sync results |

The default merge conflict mode is **Agent** -- `jmp merge` first auto-merges files with non-overlapping line changes via diff3, then invokes the preferred agent to resolve any remaining true conflicts. Use `--manual`, `--theirs`, or `--ours` to override the conflict resolution strategy (auto-merge still applies in all modes).

## Packaged workflow skill

This repository ships a small agent workflow skill at
`.agents/skills/jmp-agent-workflow/SKILL.md`. It teaches an agent to call
`jmp` directly for task start/status/finish, snapshots, drift checks, events,
and Git export while leaving code search, edits, and tests to the agent.
