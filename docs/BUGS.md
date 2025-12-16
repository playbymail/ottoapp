# Known Bugs

## Turn Reports

### Roads Text – Legacy Behavior and Limitations

The original game report generator produces the “roads” portion of the TERRAIN text from a 6-character road mask (one character per direction: N, NE, SE, S, SW, NW). For each flagged direction, it appends a road description to the TERRAIN string.

The generator attempts to avoid repeating the road type name (Rune, Dirt, Stone) by appending just the direction (e.g. “ NE”, “ SW”) if that road type has already appeared anywhere earlier in the TERRAIN string. This deduplication is applied per road type over the entire string, not per group.

When more than one road type exists in the same hex, this can produce text where later directions for an earlier road type appear after the name of a different road type. For example, a hex whose true roads are:

* Rune Road: N and SW
* Dirt Road: NE

will be printed by the legacy generator as:

```text
Rune Road N Dirt Road NE SW,
```

A straightforward reading of this text suggests:

* Rune Road: N
* Dirt Road: NE and SW

which does **not** match the actual game state. This is a known quirk of the original generator and is part of the historical report format; it cannot be corrected retroactively from the text alone. There is no deterministic way to recover the true per-type road layout from this output without consulting the game master or original data.

Our tools therefore:

* **Parse and accept the legacy text as-is.**
* **Display roads according to the literal grouping implied by the text**, not by guessing the intended game state.

If you need the true road types and directions for a hex, please contact the game master and correct your own records manually.
