# Unofficial TribeNet Turn Report Grammar

```ebnf
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


CLAN_ID      = "0" DIGIT DIGIT DIGIT
COORD        = LETTER LETTER DIGIT DIGIT DIGIT DIGIT
COORD_NA      = "N/A"
DD_MM_YYYY   = DIGIT DIGIT? "/" DIGIT DIGIT? "/" DIGIT DIGIT DIGIT DIGIT
DIGIT        = [0-9]
EOF          = !.
EOL          = [\n]
LETTER       = [A-Z]
SCOUT_ID     = [1-8]
TURN_NO      = "#" DIGIT+
TURN_YYYYMM  = DIGIT DIGIT DIGIT DIGIT? "-" DIGIT DIGIT?
