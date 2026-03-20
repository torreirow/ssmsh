## Why

In `aws-sdk-go v1.50.x` is het `Value`-veld van `ssm.Parameter` gemarkeerd als `sensitive:"true"`. De `String()`-methode van het SDK-struct (aangeroepen via `%+v` in `printResult`) redact daardoor alle parameterwaarden — ook gewone String-parameters. Daarnaast bevatte `validateTier` een case-vergelijkingsbug waardoor `tier=advanced` en `tier=standard` altijd faalden, gaf `validateOverwrite` een ruwe Go-fout terug, en produceerde `validateName` een ongeldig dubbel-slash-pad wanneer vanuit de root gewerkt werd.

## What Changes

- **Output**: `printResult` gebruikt voortaan `json.NewEncoder` met `SetEscapeHTML(false)` in plaats van `%+v`, waardoor AWS SDK sensitive-redactie wordt omzeild en `<` / `>` niet worden geëscaped.
- **SecureString masking**: `commands/get.go` vervangt de `Value` van SecureString-parameters door `<sensitive>` wanneer `decrypt=false`, en toont een waarschuwing wanneer `decrypt=true` maar de API toch `<encrypted>` retourneert (IAM/KMS-rechtenprobleem).
- **`validateTier` bugfix**: case-vergelijking gecorrigeerd van `strings.ToLower(s) == "Standard"` naar `strings.EqualFold(s, StandardTier)` — `tier=advanced/standard/ADVANCED` werken nu alle drie.
- **`validateOverwrite` foutmelding**: ruwe `strconv.ParseBool`-fout vervangen door `"overwrite must be true or false"`.
- **`validateName` padfix**: `path.Clean()` toegevoegd zodat werken vanuit root geen `//parameternaam` meer produceert.
- **Validatiefout ruis**: `shell.Println(putUsage)` verwijderd uit `validate()` — elke fout toont nu alleen de relevante foutmelding, niet de volledige usage-tekst.
- **Tests**: drie nieuwe testcases in `parameterstore_test.go` voor decryptie-gedrag en `WithDecryption`-doorgave.

## Capabilities

### New Capabilities

*(geen)*

### Modified Capabilities

- `parameter-crud`: gedrag van `get` voor SecureString-parameters gecorrigeerd; `put`-validatie robuuster gemaakt.
- `output-formatting`: tekst-output gebruikt nu JSON-serialisatie om AWS SDK sensitive-redactie te omzeilen.

## Impact

- `commands/commands.go`: `printResult` en `printJSON` herschreven; `bytes` package toegevoegd.
- `commands/get.go`: sensitive masking en `<encrypted>`-waarschuwing toegevoegd; `aws` package import toegevoegd.
- `commands/put.go`: `validateTier`, `validateOverwrite`, `validateName` gecorrigeerd; `validate()` print geen usage meer bij fouten; `path` package toegevoegd.
- `parameterstore/parameterstore_test.go`: nieuwe mock `mockedSSMCapture` + drie testfuncties.
- Geen wijzigingen in `parameterstore/parameterstore.go` — `WithDecryption` was al correct.
- Geen breaking changes voor gebruikers; output-formaat wijzigt van `%+v` naar JSON (ook zonder `output=json` in config).

## Non-goals

- Geen migratie naar AWS SDK for Go v2.
- Geen wijziging aan de `decrypt`-command interface.
- Geen fix voor `ls` (toont geen waarden).
- Geen `output`-runtime-command om het formaat te wisselen zonder config-aanpassing.
