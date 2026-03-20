# Spec: Shell Navigation

## Purpose

Provides filesystem-like navigation within the SSM Parameter Store hierarchy. The shell maintains a current working directory (cwd) and resolves relative paths against it.

## Current behaviour

### State

`ParameterStore.Cwd` holds the current path (default: `/`). The shell prompt reflects the cwd, e.g. `/dev/app>`.

### Commands

**`cd <path>`**
- Changes cwd to the given absolute or relative path.
- Relative paths are resolved via `fqp(path, cwd)` → `path.Clean()`.
- Shorthand `..` is supported (delegated to `path.Clean`).
- Errors if the resolved path does not exist in SSM (no parameters under it).

**`ls [-r|R] [path ...]`**
- Lists parameters at one or more paths (defaults to cwd if no args).
- Without `-r`: shows only top-level entries; sub-paths are shown with a trailing `/`.
- With `-r`: shows all parameters recursively.
- Internally always calls `GetParametersByPath` with `Recursive=true`, then culls results for the non-recursive case.
- Supports SIGINT (^C) to interrupt long-running listings.
- Results are sorted alphabetically before display.
- Cross-region paths supported via `region:/path` syntax.

### Path resolution (`fqp`)

```
fqp(path, cwd):
  if path starts with "/" → use as-is
  else → cwd + "/" + path
  → path.Clean()
```

## Constraints

- SSM API returns max 10 results per page; listing large hierarchies requires many paginated calls and can be slow.
- `ls` with no args on `/` at a large account may take a significant amount of time.
- `cd` to a path with no parameters returns "No such path" even if the path is syntactically valid.

## Known gaps

- No autocomplete for paths (tab completion not implemented).
- No globbing or regex filtering in `ls`.
- No `-l` (long format) to show metadata alongside names.
- No in-memory cache to speed up repeated listings.
