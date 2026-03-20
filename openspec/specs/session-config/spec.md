# Spec: Session Configuration

## Purpose

Manage the active AWS context (profile, region, KMS key) and shell behaviour (decryption, overwrite defaults) both at startup via config file and at runtime via commands.

## Config file

Default path: `~/.ssmshrc`. Override with `-config <path>` flag. INI format using `gopkg.in/gcfg.v1`.

```ini
[default]
type=SecureString      # String | StringList | SecureString
overwrite=true         # bool
decrypt=true           # bool
profile=my-profile     # AWS profile name
region=us-east-1       # AWS region
key=<kms-key-id>       # KMS key ID or ARN or alias
output=json            # text (default) | json
```

Config file is optional. If absent, all fields default to zero values (decrypt=false, overwrite=false, type=SecureString, region from env/profile).

## Priority order

| Setting | Priority 1 | Priority 2 | Priority 3 |
|---------|-----------|-----------|-----------|
| region | `AWS_REGION` env var | `.ssmshrc` | AWS profile config |
| profile | `AWS_PROFILE` env var | `.ssmshrc` | `"default"` |
| key | `.ssmshrc` | — | — |
| type | `.ssmshrc` | — | `"SecureString"` |

## Runtime commands

### `profile [name]`
- Without args: prints current profile name.
- With args: switches to named profile, re-initialises AWS session for current region.

### `region [name]`
- Without args: prints current region.
- With args: changes active region. Initialises a new SSM client for that region (added to clients map).

### `key [key-id]`
- Without args: prints current KMS key.
- With args: sets the KMS key used for new `put` operations.

### `decrypt [true|false]`
- Without args: toggles current decrypt state and prints result.
- With explicit `true`/`false`: sets state directly.
- Affects `get`, `ls`, `history`, and `cp` (cp forces decrypt=true internally).

## AWS client management

`ParameterStore.Clients` is a `map[string]ssmiface.SSMAPI` keyed by region string. Clients are lazily initialised:
- One client created at startup for the default region.
- Additional clients created on first use of a new region (via `parsePath` → `ps.InitClient`).
- `profile` change does NOT re-initialise existing region clients — only the current region's client is refreshed.

## Constraints

- No `overwrite` runtime command; overwrite can only be set in config or per-`put` via the `overwrite=` field.
- Profile switch does not invalidate cached clients for other regions (potential staleness if profile change affects cross-region permissions).
- KMS key is global across all regions in the session; a key from one region will fail in another.

## Known gaps

- No `type` runtime command (default parameter type can only be changed in config file).
- No credential validation after a `profile` switch (errors only surface on next API call).
- No support for AWS SSO / IAM Identity Center profiles beyond what the SDK handles natively.
- No `output` runtime command (output format only configurable in config file).
