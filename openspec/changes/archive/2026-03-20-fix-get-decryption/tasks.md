## 1. Diagnose

- [x] 1.1 Controleer `go.mod` op de versie van `aws-sdk-go` — vastgesteld: `v1.50.16`
- [x] 1.2 Onderzoek root cause van `<sensitive>`-output — vastgesteld: `ssm.Parameter.Value` is `sensitive:"true"` in SDK; `%+v` roept `String()` aan via `awsutil.Prettify`; `WithDecryption` in `parameterstore.go` was correct
- [x] 1.3 Verifieer of `<encrypted>` een API-gedrag of code-bug is — vastgesteld: AWS retourneert `<encrypted>` als string bij ontbrekende KMS-rechten, geen SDK-serialisatiebug

## 2. Output-fix (`commands/commands.go`)

- [x] 2.1 Vervang `shell.Printf("%+v\n", result)` door `printJSON(result)` in `printResult` om SDK-redactie te omzeilen
- [x] 2.2 Herschrijf `printJSON` met `json.NewEncoder` + `SetEscapeHTML(false)` zodat `<sensitive>` niet als `\u003csensitive\u003e` verschijnt; voeg `bytes` package toe

## 3. SecureString masking en waarschuwing (`commands/get.go`)

- [x] 3.1 Vervang `Value` van SecureString-parameters door `"<sensitive>"` wanneer `ps.Decrypt=false`
- [x] 3.2 Toon waarschuwing wanneer `ps.Decrypt=true` maar API retourneert `"<encrypted>"` (IAM/KMS-rechtenkwestie)
- [x] 3.3 Voeg `aws` package import toe

## 4. `put`-validatie fixes (`commands/put.go`)

- [x] 4.1 Fix `validateTier`: vervang `strings.ToLower(s) == StandardTier` door `strings.EqualFold(s, StandardTier/AdvancedTier)`; vervang `strings.Title` door directe constante toewijzing
- [x] 4.2 Fix `validateOverwrite`: retourneer `errors.New("overwrite must be true or false")` i.p.v. ruwe `strconv.ParseBool`-fout
- [x] 4.3 Fix `validateName`: wikkel resulterende pad in `path.Clean()` om dubbele slashes te voorkomen vanuit root; voeg `path` package toe
- [x] 4.4 Verwijder `shell.Println(putUsage)` uit `validate()` — foutmeldingen zijn nu beknopt

## 5. Tests (`parameterstore/parameterstore_test.go`)

- [x] 5.1 Voeg `mockedSSMCapture` struct toe met pointer-receiver die `in.WithDecryption` vastlegt
- [x] 5.2 `TestGetDecryptTrueReturnsPlaintext`: `Get()` met `Decrypt=true`, mock retourneert plaintext → verifieer plaintext in resultaat
- [x] 5.3 `TestGetDecryptTrueEncryptedPassthrough`: `Get()` met `Decrypt=true`, mock retourneert `<encrypted>` → verifieer pass-through (waarschuwing is verantwoordelijkheid van `commands/get.go`)
- [x] 5.4 `TestGetWithDecryptionFlagPassthrough`: verifieert dat `WithDecryption` correct wordt doorgegeven voor zowel `true` als `false`

## 6. Validatie

- [x] 6.1 `go test ./...` — alle tests slagen
- [x] 6.2 `go build ./...` — bouwt zonder fouten
- [x] 6.3 Handmatig getest met `AWS_PROFILE=TN-Production`: String-parameter toont plaintext, SecureString met `decrypt=false` toont `<sensitive>`, SecureString met `decrypt=true` toont plaintext, `put` vanuit root met naam zonder slash werkt correct
