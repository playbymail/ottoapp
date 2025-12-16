# Turn Report Structure

## Unit Location Line

```text
unitId "," SPACE UnitName? "," SPACE currentHex "," SPACE "(" previousHex ")" EOL

unitId      ::= "Tribe"    SPACE TribeId
              | "Courier"  SPACE CourierId
              | "Element"  SPACE ElementId
              | "Fleet"    SPACE FleetId
              | "Garrison" SPACE GarrisonId
currentHex  ::= "Current Hex ="  SPACE location
previousHex ::= "Previous Hex =" SPACE location
location    ::= ((Grid | "##") SPACE RowColumn)
              | "N/A"
```

Examples:

```text
Tribe 0987, , Current Hex = QQ 0203, (Previous Hex = QQ 0101)
Tribe 1987, , Current Hex = ## 0201, (Previous Hex = N/A)
Courier 2987c1, , Current Hex = WX 0220, (Previous Hex = WX 1001)
Element 3987e1, , Current Hex = PR 1801, (Previous Hex = PQ 1820)
Fleet 4987f1, , Current Hex = QS 1201, (Previous Hex = SQ 2001)
Garrison 5987g1, , Current Hex = QQ 0303, (Previous Hex = QQ 0303)
```

## Turn Information Line

```text
currentTurn ( TAB nextTurn )?

currentTurn ::= "Current Turn"  SPACE YearDashMonth SPACE "(#" TurnNo ")," SPACE Season "," SPACE Weather
nextTurn    ::= TAB "Next Turn" SPACE YearDashMonth SPACE "(#" TurnNo ")," SPACE reportDate
reportDate  ::= DaySlashMonthSlashYear
```

Examples:

```text
Current Turn 904-01 (#49), Spring, FINE	Next Turn 904-02 (#50), 28/09/2025
Current Turn 904-01 (#49), Spring, FINE
```

## Tribe Follows Line

```text
"Tribe Follows" SPACE UnitId
```

Example:

```text
Tribe Follows 0987e1
```
## Tribe Goes To Line

```text
"Tribe Goes to" SPACE Grid SPACE RowColumn
```

Example:

```text
Tribe Goes to QA 1211
```

## Tribe Movement Line


## Tribe Movement Failed Line

```text
"Tribe Movement: Not enough animals to pull wagons. Movement is not possible."
```

Examples:

```text
"Tribe Movement: Not enough animals to pull wagons. Movement is not possible."
```

```text
NO_MOVEMENT_REASON = "Can't Move on Lake to " & Direction
NO_MOVEMENT_REASON = "Can't Move on Ocean to " & Direction
NO_MOVEMENT_REASON = "Cannot Move Wagons into Jungle Hill to " & Direction & " of HEX"
NO_MOVEMENT_REASON = "Cannot Move Wagons into Mountains  to " & Direction & " of HEX"
NO_MOVEMENT_REASON = "Cannot Move Wagons into Snowy hills to " & Direction & " of HEX"
NO_MOVEMENT_REASON = "Cannot Move Wagons into Swamp  to " & Direction & " of HEX"
NO_MOVEMENT_REASON = "Horses not allowed into MANGROVE Swamp to " & Direction & " of HEX"
NO_MOVEMENT_REASON = "Insufficient capacity to carry "
NO_MOVEMENT_REASON = "Insufficient capacity to carry"
NO_MOVEMENT_REASON = "NO DIRECTION"
NO_MOVEMENT_REASON = "No Ford on River to " & Direction & " of HEX"
NO_MOVEMENT_REASON = "No River Adjacent to Hex to " & Direction
NO_MOVEMENT_REASON = "Not enough animals to pull wagons"
NO_MOVEMENT_REASON = "Not enough M.P's to move to " & Direction & " into " & NEW_TERRAIN
NO_MOVEMENT_REASON = "Not enough M.P's"
NO_MOVEMENT_REASON = NO_MOVEMENT_REASON & " of HEX "
NO_MOVEMENT_REASON = NO_MOVEMENT_REASON & " of HEX"
```

## Fleet Movement Line

## Scouting Line

## Unit Status Line

```text
UnitId " Status: "
  (
    Terrain ","
    ( SPACE SpecialHex "," )?
    ( SPACE City       "," )?
    ( SPACE Ore        "," )*
    ( SPACE "Lcm" SPACE Direction ( "," SPACE Direction )*  "," )?
    ( SPACE "Lsm" SPACE Direction ( "," SPACE Direction )*  "," )?
    ( SPACE "Hsm" SPACE Direction ( "," SPACE Direction )*  "," )?
    ( SPACE "LVm" SPACE Direction ( "," SPACE Direction )*  "," )?
    ( SPACE "LJm" SPACE Direction ( "," SPACE Direction )*  "," )?
    ( SPACE "O"   SPACE Direction ( "," SPACE Direction )*  "," )?
    ( SPACE "L"   SPACE Direction ( "," SPACE Direction )*  "," )?
    ( ",Pass"     SPACE Direction (     SPACE Direction )*      )?
    ( ",River"    SPACE Direction (     SPACE Direction )*      )?
    ( ",Ford"     SPACE Direction (     SPACE Direction )*      )?
    ( ",Canal"    SPACE Direction (     SPACE Direction )*      )?
    ((SPACE ("Rune"|"Dirt"|"Stone") SPACE "Road" (SPACE Direction))+ ",")?
    ( SPACE "Quarry Hex," )?
    ( "," SPACE "Whaling Area (Improved Whaling)" )?
  )
  SPACE ( "," SPACE UnitId )*
  EOL
```

Notes:

1. Directions for mountains, oceans, and lakes are ordered NE SE SW NW N S.
2. Directions for passes, rivers, fords, canals, and roads are ordered N NE SE S SW NW.

## Comma Quirks

```vba
Function CleanReportString(ByVal txt As String) As String
    txt = Replace( txt , ",,"  , "," ) ' replace comma       comma with comma
    txt = Replace( txt , " ,"  , "," ) ' replace       space comma with comma
    txt = Replace( txt , ", ," , "," ) ' replace comma space comma with comma
    CleanReportString = txt
End Function
```

## Missing
STILL EMPTY HALT
