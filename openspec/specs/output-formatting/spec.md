# Spec: Output Formatting

## Purpose

Control how parameter data is presented to the user after `get`, `ls`, `history`, and `put` operations.

## Current behaviour

### Output modes

Configured via `.ssmshrc` `output` field. No runtime command to change it.

| Mode | Config value | How set |
|------|-------------|---------|
| JSON (default én expliciete modus) | omitted of `json` | `json.NewEncoder` met `SetEscapeHTML(false)` |

Beide modi gebruiken JSON-serialisatie. `%+v` is verwijderd omdat de AWS SDK `sensitive`-tag op `ssm.Parameter.Value` alle waarden redact in `aws-sdk-go v1.50.x`.

### Per-command output

**`ls`**
- Always text. Prints one name per line.
- Sub-paths shown with trailing `/` (e.g. `app/`).
- Results sorted alphabetically.
- No metadata (type, version, modified date).

**`get`**
- JSON via `json.NewEncoder` met `SetEscapeHTML(false)`.
- Veld `Value` van SecureString-parameters wordt vervangen door `"<sensitive>"` wanneer `decrypt=false`.
- Velden komen overeen met AWS SDK Go struct namen (PascalCase).

**`put`**
- Always text. Prints `Put <name> version <n>` on success.
- Errors printed via `shell.Println`.

**`history`**
- JSON via `json.NewEncoder` met `SetEscapeHTML(false)`.
- Velden komen overeen met AWS SDK Go struct namen (PascalCase).

**`cp`, `mv`, `rm`**
- No output on success (silent).
- Errors printed via `shell.Println`.

### JSON field names (examples for `get`)

```json
{
    "ARN": "arn:aws:ssm:...",
    "DataType": "text",
    "LastModifiedDate": "2024-01-01T00:00:00Z",
    "Name": "/dev/db/username",
    "Selector": null,
    "SourceResult": null,
    "Type": "SecureString",
    "Value": "<sensitive>",
    "Version": 1
}
```

Field names come directly from the AWS SDK Go struct tags — they are not customised.

## Constraints

- Output mode is global (one setting for all commands).
- `ls` output is not affected by the `output` setting (always plain text).
- JSON uses AWS SDK struct names (PascalCase), not idiomatic JSON (camelCase).
- No table/column-aligned output.
- Tekens zoals `<` en `>` in waarden worden niet geëscaped (`SetEscapeHTML(false)`).

## Known gaps

- No `--output` flag per command to override global setting at runtime.
- No table output format (e.g. aligned columns for name/type/version).
- No YAML output format.
- No CSV/TSV output for scripting pipelines.
- No way to select specific fields (e.g. `get /foo --field Value`).
- `ls` does not show metadata even in JSON mode.
- No `--quiet` flag to suppress version output from `put`.
