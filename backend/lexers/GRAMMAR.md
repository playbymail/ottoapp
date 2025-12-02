# Unofficial TribeNet Turn Report Grammar

## Lemon like

```lemon
turnReport ::= clanSection EOF .

clanSection ::= clanLocation .

clanLocation ::= CLAN_ID COMMA note COMMA CURRENT_HEX EQ coords COMMA LPAREN PREVIOUS_HEX EQ coords RPAREN EOL .

coords := COORD .
coords ::= COORD_NA .
coords ::= COORD_OBSCURED .
```

### Terminals

These terminals are defined using Go style regular expressions.

```go
type TokenKind_t string
var (
    // EOF token is emitted by the scanner when input is exhausted
    EOF             TokenKind_t = "EOF"
    CLAN_ID         TokenKind_t = "CLAN_ID"        // regexp.MustCompile(`Tribe 0\d{3}`)
    COORD           TokenKind_t = "COORD"          // regexp.MustCompile(`[A-Z]{2} \d{4}`)
    COORD_NA        TokenKind_t = "COORD_NA"       // regexp.MustCompile(`N/A`)
    COORD_OBSCURED  TokenKind_t = "COORD_OBSCURED" // regexp.MustCompile(`## \d{4}`)
    COURIER_ID      TokenKind_t = "COURIER_ID"     // regexp.MustCompile(`\d{4}c[1-9]`)
    DATE            TokenKind_t = "DATE"           // regexp.MustCompile(`\d{1,2}/\d{1,2}/\d{4}`)
    ELEMENT_ID      TokenKind_t = "ELEMENT_ID"     // regexp.MustCompile(`\d{4}e[1-9]`)
    EOL             TokenKind_t = "EOL"            // regexp.MustCompile(`\n`)
    FLEET_ID        TokenKind_t = "FLEET_ID"       // regexp.MustCompile(`\d{4}f[1-9]`)
    GARRISON_ID     TokenKind_t = "GARRISON_ID"    // regexp.MustCompile(`\d{4}g[1-9]`)
    SCOUT_ID        TokenKind_t = "SCOUT_ID"       // regexp.MustCompile(`Scout [1-8]`)
    TRIBE_ID        TokenKind_t = "TRIBE_ID"       // regexp.MustCompile(`[1-9]\d{3}`)
    TURN_NO         TokenKind_t = "TURN_NO"        // regexp.MustCompile(`#\d{1,5}`)
    TURN_YYYYMM     TokenKind_t = "TURN_YYYYMM"    // regexp.MustCompile(`\d{3,4}-\d{1,2}`)
    UNIT_ID         TokenKind_t = "UNIT_ID"        // regexp.MustCompile(`\d{4}([cefg][1-9])?`)
)

```


