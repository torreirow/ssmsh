## Context

Drie losstaande problemen zijn tijdens de implementatie vastgesteld:

**1. AWS SDK sensitive-redactie (root cause `<sensitive>`)**
In `aws-sdk-go v1.50.x` is `ssm.Parameter.Value` gemarkeerd als `sensitive:"true"`. De `String()`-methode van het struct (gegenereerd via `awsutil.Prettify`) redact dit veld automatisch. `printResult` gebruikte `%+v` als format-verb, wat `String()` aanroept — waardoor alle parameterwaarden, ook gewone String-parameters, als `<sensitive>` verschenen. `WithDecryption` in `parameterstore.go` was correct; het probleem zat uitsluitend in de output-laag.

**2. `put`-validatie bugs**
Drie problemen in `commands/put.go`:
- `validateTier`: vergeleek `strings.ToLower(s)` met `"Standard"`/`"Advanced"` (hoofdletters) — geen enkele waarde haalde ooit de check.
- `validateOverwrite`: retourneerde de ruwe `strconv.ParseBool`-fout in plaats van een gebruiksvriendelijke melding.
- `validateName`: bouwde pad als `ps.Cwd + "/" + name` zonder `path.Clean()`, wat bij cwd=`/` een dubbele slash produceerde (`//parameternaam`) die AWS weigert.
- `validate()`: printte bij elke validatiefout de volledige usage-tekst, wat ruis veroorzaakte.

## Goals / Non-Goals

**Goals:**
- Correcte weergave van parameterwaarden in text-output (omzeilen SDK-redactie).
- Expliciete `<sensitive>`-masking voor SecureString met `decrypt=false`.
- Waarschuwing bij `<encrypted>` retour (IAM/KMS-rechtenkwestie).
- Correcte `put`-validatie voor tier, overwrite, en namen vanuit root.
- Testdekking voor decryptie-gedrag en `WithDecryption`-doorgave.

**Non-Goals:**
- Migratie naar AWS SDK for Go v2.
- Wijzigingen aan de `decrypt`-command interface.
- Runtime `output`-command om formaat te wisselen.

## Decisions

### 1. JSON-encoder als universele output-methode

**Beslissing:** `printResult` gebruikt voor beide output-modi (text én json) `json.NewEncoder` met `SetEscapeHTML(false)`. `%+v` is verwijderd.

**Alternatieven overwogen:**
- Velden handmatig per struct type afdrukken → te veel koppeling met AWS SDK-structs.
- `awsutil.Prettify` direct aanroepen → redact nog steeds `sensitive`-velden, lost het probleem niet op.
- `%v` i.p.v. `%+v` → zelfde probleem, roept ook `String()` aan.

**Rationale:** `json.MarshalIndent` serialiseert via JSON-tags en omzeilt de `Stringer`-interface volledig. `SetEscapeHTML(false)` voorkomt dat `<sensitive>` als `\u003csensitive\u003e` verschijnt.

### 2. SecureString masking in `commands/get.go`

**Beslissing:** Na ontvangst van de API-respons vervangt `get` de `Value` van SecureString-parameters door `"<sensitive>"` wanneer `ps.Decrypt=false`. Bij `ps.Decrypt=true` en een geretourneerde waarde `"<encrypted>"` wordt een waarschuwing getoond.

**Rationale:** De presentatie-laag is de juiste plek voor dit gedrag; `parameterstore.Get()` blijft puur een data-doorgeefluik zonder output-logica.

### 3. `validateTier` fix via `strings.EqualFold`

**Beslissing:** Vervang de kapotte `strings.ToLower(s) == StandardTier`-vergelijking door `strings.EqualFold(s, StandardTier)`. `strings.Title` (deprecated) vervangen door directe constante toewijzing.

**Rationale:** `strings.EqualFold` is de idiomatische Go-manier voor case-insensitieve vergelijking.

### 4. `validateName` padnormalisatie via `path.Clean`

**Beslissing:** Wikkel het resulterende pad in `path.Clean()`, consistent met `fqp()` in `parameterstore.go`.

**Rationale:** `fqp()` is niet geëxporteerd; `path.Clean()` direct aanroepen in `validateName` volgt hetzelfde patroon zonder koppeling aan het `parameterstore`-package.

### 5. Testdekking via `mockedSSMCapture`

**Beslissing:** Nieuwe mock-struct `mockedSSMCapture` met pointer-receiver die `in.WithDecryption` vastlegt, naast de bestaande value-type `mockedSSM`. Drie nieuwe testfuncties in `parameterstore_test.go`.

**Rationale:** De bestaande `mockedSSM` hardcodde `WithDecryption: aws.Bool(true)` en kon de doorgegeven vlag niet valideren. Een aparte mock voorkomt breaking changes in bestaande tests.

## Risks / Trade-offs

- **[Trade-off] Output-formaat gewijzigd** → Text-output is nu altijd JSON, ook zonder `output=json` in config. Bestaande scripts die de `%+v`-opmaak parseerden kunnen breken. Mitigatie: JSON is een stabielere en beter parseerbare opmaak.
- **[Trade-off] Geen integratie-test tegen echte AWS** → Mock-tests valideren logica maar niet het API-gedrag. Acceptabel omdat een echte AWS-omgeving buiten de CI-scope valt.
- **[Risico] `<encrypted>` als waarschuwing, niet als fout** → Als IAM-rechten ontbreken, gaat de put of get door maar de waarde is onbruikbaar. Mitigatie: waarschuwing is zichtbaar in de output vóór het resultaat.

## Open Questions

*(geen — alle vragen zijn beantwoord tijdens implementatie)*
