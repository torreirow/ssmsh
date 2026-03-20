# Spec: Parameter Operations

## Purpose

Higher-level operations on parameters: copy, move, history retrieval, and named parameter policies.

## Current behaviour

### `cp [-r|R] <src> <dst>`

Copies a parameter or path hierarchy to a new location.

**Decision matrix:**

| src type | dst type | behaviour |
|----------|----------|-----------|
| parameter | non-existent | copy param to dst name |
| parameter | existing path | copy param into path (preserving param name) |
| path | parameter | error |
| path | non-existent path | copy all params under src to dst (new path) |
| path | existing path | copy all params under src into dst/src-basename |

- Forces `Decrypt=true` temporarily during copy so that SecureString values are readable for re-encryption at destination.
- Uses `GetParameterHistory` to retrieve the latest version's metadata (type, key, description, allowed pattern) for faithful reproduction.
- Supports cross-region: `cp -r us-east-1:/dev us-west-2:/dev`.
- Pagination handled for path-to-path copies (SSM returns max 10 per page).

### `mv <src> <dst>`

Move = `cp` (recursive=true) followed by `rm` (recursive=true) on src.

- No atomicity guarantee: if the `rm` fails after a successful `cp`, the source is left in place.
- Cross-region moves are supported (copies to target region, deletes from source).

### `history <param>`

Retrieves full version history of a single parameter.

- Calls `GetParameterHistory` with pagination.
- Respects `ParameterStore.Decrypt` flag.
- Output format: text or JSON (controlled by config).
- History entries include: name, type, value, key, description, labels, tier, policies, last modified date/user, version number.

### `policy <name> <policy-expression> [policy-expression ...]`

Defines a named in-memory policy object for use with `put policies=[name1,name2]`.

**Supported policy types:**

| Expression | Description |
|------------|-------------|
| `Expiration(Timestamp=<ISO8601>)` | Auto-delete parameter at timestamp |
| `ExpirationNotification(Before=<n>,Unit=<days\|hours>)` | Notify before expiration |
| `NoChangeNotification(After=<n>,Unit=<days\|hours>)` | Notify if unchanged for duration |

- Policies are stored in a package-level `policies` map (in-memory only, not persisted).
- Using any policy in `put` automatically sets `tier=Advanced`.
- Multiple named policies can be combined in a single `put`.

## Constraints

- `mv` is not atomic: a crash or error between copy and delete leaves duplicate data.
- `cp` cross-region requires both regions to use compatible KMS keys (or AWS managed key).
- Named policies are session-local; they are lost when the shell exits.
- `history` only works on a single parameter at a time.

## Known gaps

- No cross-account copy (only cross-region within same credentials).
- No `--dry-run` mode for `cp`/`mv` to preview what would be affected.
- No `tag` command to manage SSM parameter tags.
- Policies are not persisted between sessions.
