# Unofficial TribeNet Turn Report Grammar

## Lemon like

```lemon

turnReport ::= clanSection EOF .

clanSection ::= clanLocation .

clanLocation ::= CLAN_ID COMMA note COMMA SP CURRENT_HEX SP EQ coords COMMA SP LPAREN PREVIOUS_HEX SP EQ coords RPAREN EOL .

coords := COORD .
coords ::= COORD_NA .
coords ::= COORD_OBSCURED .

TurnReport <- ClanSection EOF {
    ...
}

ClanSection <- ClanLocation CurrentTurnLong ScoutLines* {
    ...
}
 
ClanLocation <- "Tribe " ClanId ", , Current Hex = " Coord ", (Previous Hex = " Coord "(" EOL {
    ...
}

Coord <- (COORD / COORD_NA) {
    ...
}


CurrentTurnLong <- "Current Turn " TURN_YYYYMM " (" TURN_NO "), " Season ", " Weather " Next Turn " TURN_YYYYMM " (" TURN_NO "), " DD_MM_YYYY EOL {
    ...
}

Season <- "Winter" {
    ...
}

Weather <= "FINE" {
    ...
}

ScoutLine <- "Scout " SCOUT_ID ":Scout " ScoutMove (BACKSLASH ScoutMove)* ScoutPatrolled? EOL {
    ...
}

ScoutMove <- (ScoutMoved / ScoutFailed) {
    ...
}

ScoutMoved <- DIRECTION "-" TERRAIN ScoutMoveResults {
    ...
}

ScoutFailed <- ",Can't Move on Lake to " d:Direction " of HEX" {
   ...
} / ",Can't Move on Ocean to " d:Direction " of HEX" {
  ...
} / ",No Ford on River to " d:Direction " of HEX" {
  ...
} / ",Not enough M.P's to move to " d:Direction " into " t:TerrainLong {
  ...
}

Direction <- ("NE" / "NW" / "SE" / "SW" / "N" / "S") {
  ...
}

Terrain <- ("BF" / "GH" / "JH" / "PR") {
  ...
}

TerrainLong <- ("BRUSH FLAT" / "GRASSY HILLS" / "PRAIRIE") {
  ...
}
```


### Terminals

These terminals are defined using Go style regular expressions.

```go
// Note: EOF token is emitted by the scanner when input is exhausted
type TokenType string
var (
    EOF             TokenType = "EOF"
    CLAN_ID         TokenType = "CLAN_ID"        // regexp.MustCompile(`Tribe 0\d{3}`)
    COORD           TokenType = "COORD"          // regexp.MustCompile(`[A-Z]{2} \d{4}`)
    COORD_NA        TokenType = "COORD_NA"       // regexp.MustCompile(`N/A`)
    COORD_OBSCURED  TokenType = "COORD_OBSCURED" // regexp.MustCompile(`## \d{4}`)
    COURIER_ID      TokenType = "COURIER_ID"     // regexp.MustCompile(`\d{4}c[1-9]`)
    DATE            TokenType = "DATE"           // regexp.MustCompile(`\d{1,2}/\d{1,2}/\d{4}`)
    ELEMENT_ID      TokenType = "ELEMENT_ID"     // regexp.MustCompile(`\d{4}e[1-9]`)
    EOL             TokenType = "EOL"            // regexp.MustCompile(`\n`)
    FLEET_ID        TokenType = "FLEET_ID"       // regexp.MustCompile(`\d{4}f[1-9]`)
    GARRISON_ID     TokenType = "GARRISON_ID"    // regexp.MustCompile(`\d{4}g[1-9]`)
    SCOUT_ID        TokenType = "SCOUT_ID"       // regexp.MustCompile(`Scout [1-8]`)
    TRIBE_ID        TokenType = "TRIBE_ID"       // regexp.MustCompile(`[1-9]\d{3}`)
    TURN_NO         TokenType = "TURN_NO"        // regexp.MustCompile(`#\d{1,5}`)
    TURN_YYYYMM     TokenType = "TURN_YYYYMM"    // regexp.MustCompile(`\d{3,4}-\d{1,2}`)
    UNIT_ID         TokenType = "UNIT_ID"        // regexp.MustCompile(`\d{4}([cefg][1-9])?`)
)

```


