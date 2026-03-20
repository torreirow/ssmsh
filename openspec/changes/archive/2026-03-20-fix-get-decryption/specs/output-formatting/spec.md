## MODIFIED Requirements

### Requirement: text-output gebruikt JSON-serialisatie
De standaard text-output (zonder `output=json` in config) MOET JSON-serialisatie gebruiken via `json.NewEncoder` met `SetEscapeHTML(false)`. Het systeem MAG NIET `%+v` gebruiken als format-verb op AWS SDK-structs, omdat dit de `String()`-methode aanroept die `sensitive`-getagde velden redact.

#### Scenario: get zonder output=json toont waarden
- **WHEN** de gebruiker `get /pad/naar/param` uitvoert zonder `output=json` in de config
- **THEN** toont de shell de parameterwaarde in JSON-formaat met de werkelijke waarde, niet `<sensitive>`

#### Scenario: speciale tekens in waarden worden niet geëscaped
- **WHEN** een parameterwaarde of veldnaam de tekens `<` of `>` bevat (bijv. `<sensitive>`)
- **THEN** toont de shell de tekens letterlijk, niet als `\u003c` of `\u003e`

## REMOVED Requirements

### Requirement: text-output gebruikt Go struct-formaat (`%+v`)
**Reason:** `%+v` roept de `String()`-methode aan van AWS SDK-structs, die via `awsutil.Prettify` alle `sensitive:"true"`-getagde velden redact. Vanaf `aws-sdk-go v1.50.x` is `ssm.Parameter.Value` als `sensitive` gemarkeerd, waardoor alle parameterwaarden als `<sensitive>` verschijnen ongeacht het parametertype.
**Migration:** Zowel de standaard text-modus als de expliciete `output=json`-modus gebruiken nu JSON-serialisatie. Er is geen gedragsverschil meer tussen de twee modi.
