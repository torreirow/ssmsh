# Spec: Parameter CRUD

## Purpose

Create, read, and delete individual parameters or entire path hierarchies in SSM Parameter Store.

## Current behaviour

### `get <param> [param ...]`

- Fetches one or more parameters by name (absolute or relative to cwd).
- Groups parameters by region before calling `GetParameters` (batch API).
- `WithDecryption` wordt ingesteld op de actuele waarde van `ParameterStore.Decrypt` op het moment van de aanroep — nooit gecached of gereset.
- SecureString-parameters met `decrypt=false`: `Value` wordt getoond als `"<sensitive>"`.
- SecureString-parameters met `decrypt=true`: `Value` bevat de plaintext-waarde.
- Als de SSM API `"<encrypted>"` retourneert terwijl `decrypt=true` is ingesteld, toont de shell een waarschuwing dat KMS-rechten mogelijk ontbreken.
- Output altijd via JSON-serialisatie (`json.NewEncoder`, `SetEscapeHTML(false)`); AWS SDK `sensitive`-tag heeft geen invloed op de weergegeven waarden.
- Cross-region: `us-east-1:/dev/db/password`.
- Returns nothing (no error) when a parameter does not exist — SSM `GetParameters` silently omits missing names from the response.

### `put` (inline or multiline)

**Inline:**
```
put name=/path/to/param value="val" type=String [options...]
```

**Multiline:**
```
put
... name=/path/to/param
... value=val
... type=String
...
```

**Fields:**

| Field | Required | Notes |
|-------|----------|-------|
| name | Yes | Absolute or relative to cwd; `path.Clean()` toegepast op resultaat |
| value | Yes | Multiline values supported |
| type | Yes | `String`, `StringList`, `SecureString` |
| description | No | |
| key | No | KMS key ARN or alias |
| pattern | No | Regex allowed pattern |
| overwrite | No | `true`/`false`; defaults to `ParameterStore.Overwrite` |
| region | No | Target region for this put |
| tier | No | `standard` of `advanced` (case-insensitief) |
| policies | No | Named policies (see policy command); forces `advanced` tier |

- Defaults for type, key, overwrite are taken from `ParameterStore` state (set via config or runtime commands).
- Validatie is veld-voor-veld via een `map[string]func` dispatch table.
- Bij een validatiefout: alleen de relevante foutmelding wordt getoond (geen usage-dump), `putParamInput` wordt gereset.
- `tier`-validatie is case-insensitief: `standard`, `Standard`, `STANDARD`, `advanced`, `Advanced`, `ADVANCED` zijn alle geldig.
- `overwrite`-validatie retourneert `"overwrite must be true or false"` bij ongeldige invoer.
- `name` zonder leading slash vanuit cwd `/` produceert `/naam` (niet `//naam`) dankzij `path.Clean()`.

### `rm [-r|R] <param> [param ...]`

- Removes one or more parameters.
- Without `-r`: parameter must be an exact parameter name, not a path.
- With `-r`: recursively deletes all parameters under a path.
- Batches deletions in groups of 10 (SSM API limit for `DeleteParameters`).
- Groups by region before deleting.
- Returns an error listing any invalid (non-existent) parameter names returned by SSM.

## Constraints

- SSM `GetParameters` is limited to 10 names per call (the SDK handles this differently from `GetParametersByPath`).
- `put` with `type=SecureString` requires a KMS key; if none is set the AWS default key (`alias/aws/ssm`) is used.
- `rm` without `-r` on a path-only entry (no direct parameter) returns an error rather than silently skipping.
- Downgraden van `tier=advanced` naar `tier=standard` op een bestaande parameter is niet mogelijk (AWS-beperking).

## Known gaps

- `get` does not warn when a requested parameter does not exist.
- No `find` / search command to locate parameters by name pattern or value.
- No glob/regex support in `get` or `rm`.
- No bulk `get` via path (use `ls -r` + multiple `get` calls as workaround).
