## MODIFIED Requirements

### Requirement: get retourneert ontsleutelde waarden bij decrypt=true
Wanneer `ParameterStore.Decrypt` `true` is, MOET de `get`-command de plaintext-waarde van SecureString-parameters retourneren. Het systeem MAG NIET de string `<encrypted>` als waarde tonen. Als de SSM-service `<encrypted>` teruggeeft omdat KMS-rechten ontbreken, MOET het systeem een waarschuwing tonen vóór het resultaat.

#### Scenario: get met decrypt=true op SecureString
- **WHEN** `ParameterStore.Decrypt` is `true` en de gebruiker voert `get /pad/naar/param` uit op een SecureString-parameter
- **THEN** toont de shell de plaintext-waarde van de parameter

#### Scenario: get met decrypt=true maar zonder KMS-rechten
- **WHEN** `ParameterStore.Decrypt` is `true` en de SSM API retourneert de letterlijke string `<encrypted>`
- **THEN** toont de shell een waarschuwing: `Warning: <naam> returned <encrypted>. Verify that your IAM role has kms:Decrypt permission for the parameter's KMS key.`

#### Scenario: get met decrypt=false op SecureString
- **WHEN** `ParameterStore.Decrypt` is `false` en de gebruiker voert `get /pad/naar/param` uit op een SecureString-parameter
- **THEN** toont de shell `"<sensitive>"` als waarde — de ruwe versleutelde blob wordt niet getoond

### Requirement: put valideert tier case-insensitief
De `put`-command MOET `tier=standard`, `tier=Standard`, `tier=STANDARD`, `tier=advanced`, `tier=Advanced` en `tier=ADVANCED` allemaal accepteren als geldige invoer.

#### Scenario: tier met kleine letters
- **WHEN** de gebruiker `put ... tier=advanced` uitvoert
- **THEN** wordt de parameter aangemaakt met tier `Advanced` en geen validatiefout

#### Scenario: ongeldige tier
- **WHEN** de gebruiker `put ... tier=premium` uitvoert
- **THEN** toont de shell `tier must be standard or advanced` en wordt de put afgebroken

### Requirement: put toont beknopte foutmeldingen
Bij een validatiefout MOET de `put`-command alleen de relevante foutmelding tonen. Het systeem MAG NIET de volledige usage-tekst afdrukken bij elke validatiefout.

#### Scenario: ongeldige overwrite-waarde
- **WHEN** de gebruiker `put ... overwrite=misschien` uitvoert
- **THEN** toont de shell alleen `overwrite must be true or false`

#### Scenario: ongeldig type
- **WHEN** de gebruiker `put ... type=Bogus` uitvoert
- **THEN** toont de shell alleen `Invalid type Bogus`

### Requirement: put lost relatieve namen op vanuit root correct op
Wanneer `ps.Cwd` gelijk is aan `/` en de gebruiker een naam opgeeft zonder leading slash, MOET de `put`-command het pad normaliseren naar `/<naam>` zonder dubbele slash.

#### Scenario: put zonder leading slash vanuit root
- **WHEN** `ps.Cwd` is `/` en de gebruiker `put name=mijnparam ...` uitvoert
- **THEN** wordt de parameter aangemaakt als `/mijnparam`, niet als `//mijnparam`

## ADDED Requirements

### Requirement: WithDecryption-vlag wordt consistent doorgegeven
De functies `Get`, `List` en `GetHistory` in `parameterstore.go` MOETEN elk de `WithDecryption`-vlag instellen op de actuele waarde van `ParameterStore.Decrypt` op het moment van de aanroep.

#### Scenario: Decrypt-toggle gevolgd door get
- **WHEN** de gebruiker `decrypt true` uitvoert en daarna `get /pad/param`
- **THEN** wordt `WithDecryption=true` meegestuurd in het `GetParametersInput`-verzoek

### Requirement: output-formaat omzeilt SDK-redactie
De `get`- en `history`-commands MOETEN parameterwaarden altijd tonen via JSON-serialisatie, zodat de AWS SDK `sensitive`-tag geen invloed heeft op de weergegeven waarden.

#### Scenario: String-parameter toont altijd plaintext
- **WHEN** de gebruiker `get /pad/naar/stringparam` uitvoert op een parameter van type `String`
- **THEN** toont de shell de werkelijke waarde, niet `<sensitive>`
