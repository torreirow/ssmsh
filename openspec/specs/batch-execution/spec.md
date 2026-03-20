# Spec: Batch Execution

## Purpose

Execute ssmsh commands non-interactively: from a file, from stdin, or as a single inline command. Enables scripting and automation workflows.

## Current behaviour

### Modes (mutually exclusive, checked in order)

```
ssmsh -file -           → stdin mode
ssmsh -file <path>      → file mode
ssmsh <cmd> [args...]   → inline mode (len(flag.Args()) > 1)
ssmsh                   → interactive mode (default)
```

### File mode (`-file <path>`)

Reads the entire file, splits on newlines, executes each line sequentially via `shell.Process`.

- Empty lines are skipped.
- Lines starting with `#` are treated as comments and skipped.
- Lines are parsed with `go-shellwords` (handles quoted strings, escape sequences).
- On first error: prints error and exits with code 1.
- No partial execution recovery — a failure stops all remaining commands.

### Stdin mode (`-file -`)

Reads all of stdin via `ioutil.ReadAll`, then processes identically to file mode.

Typical use:
```bash
cat commands.txt | ssmsh -file -
echo 'get /dev/db/password' | ssmsh -file -
```

### Inline mode

When `len(flag.Args()) > 1`, the remaining args after flags are passed directly to `shell.Process` as a single command.

```bash
ssmsh put name=/dev/app/url value="https://example.com" type=String
```

Note: `flag.Args()` starts at index 0 for the first non-flag argument. The condition `> 1` means at least two non-flag args are required (command + at least one argument). A bare `ssmsh get` with no path would not trigger inline mode — it would open interactive mode.

## Flags

| Flag | Description |
|------|-------------|
| `-config <path>` | Load config from specified file (default: `~/.ssmshrc`) |
| `-file <path\|->` | Read commands from file or stdin |
| `-version` | Print version and exit |

## Error handling

- Config read error → exit 1.
- AWS credential/session error → exit 1.
- Command parse error (shellwords) → exit 1.
- Command execution error → exit 1.
- All errors print a message to stdout (via `shell.Println`).

## Constraints

- Inline mode requires at least 2 non-flag arguments; single-word commands (e.g. `ssmsh decrypt`) open interactive mode instead.
- No `--quiet` / `--silent` flag to suppress output in batch mode.
- No exit-code propagation per command; only first error causes non-zero exit.
- `ioutil.ReadAll` (deprecated since Go 1.16) is used for stdin reading.

## Known gaps

- No `--dry-run` flag for batch files.
- No line-number reporting in error messages for file mode.
- No support for shell conditionals or variables within batch files.
- No way to continue on error (`set -e` equivalent is hardcoded on).
