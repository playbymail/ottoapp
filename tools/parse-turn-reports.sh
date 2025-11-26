#!/bin/bash

OTTOAPP_ROOT=/Users/wraith/Software/playbymail/ottoapp
OTTOMAP_ROOT=/Users/wraith/Jetbrains/tribenet/ottomap
TN0301_ROOT=/Users/wraith/Software/playbymail/ottoapp/data/tn3.1

cd ${OTTOAPP_ROOT} || exit 2
cd ${OTTOAPP_ROOT} || exit 2
cd ${TN0301_ROOT} || exit 2

cd ${OTTOAPP_ROOT} || exit 2
go build -o ${TN0301_ROOT}/bin/ottoapp ./cmd/ottoapp || exit 2
OTTOAPP_DBPATH=${OTTOAPP_ROOT}/data/alpha
[ -d "${OTTOAPP_DBPATH}" ] || {
  echo "error: bad database path '${OTTOAPP_DBPATH}"
  exit 2
}
OTTOAPP_DB=${OTTOAPP_DBPATH}/ottoapp.db
[ -f "${OTTOAPP_DB}" ] || {
  echo "error: bad database path '${OTTOAPP_DB}"
  exit 2
}

cd ${OTTOMAP_ROOT} || exit 2
go build -o ${TN0301_ROOT}/bin/ottomap || exit 2

cd ${TN0301_ROOT} || exit 2
echo " info: current ottomap version: $( bin/ottomap version )"
echo " info: current ottoapp version: $( bin/ottoapp version )"

# flatten the game data into a simple map of clanNo to handle
jq '
  def lpad(n; ch):
    tostring
    | (n - length) as $pad
    | if $pad > 0 then (ch * $pad) + . else . end;

  .players
  | to_entries
  | map(
      . as $entry
      | $entry.value.games
      | map(
          select(.id == "0301")
          | { (.clan | lpad(4; "0")): $entry.key }
        )
    )
  | add | add
' ../inputs/alpha.json > clans.json || exit 2

# verify the clan setup before uploading the map file
foundErrors=
jq -r 'to_entries | sort_by(.key) | .[] | "\(.key) \(.value)"' clans.json |
while read -r clanNo handle; do
  [ -d "${clanNo}" ] || {
    foundErrors=TRUE
    echo "error: clan ${clanNo}: missing clan folder"
    continue
  }
  [ -d "${clanNo}/data" ] || {
    foundErrors=TRUE
    echo "error: clan ${clanNo}: missing clan data folder"
    continue
  }
  [ -d "${clanNo}/data/input" ] || {
    foundErrors=TRUE
    echo "error: clan ${clanNo}: missing clan data/input folder"
    continue
  }
  [ -d "${clanNo}/data/logs" ] || {
    foundErrors=TRUE
    echo "error: clan ${clanNo}: missing clan data/logs folder"
    continue
  }
  [ -d "${clanNo}/data/output" ] || {
    foundErrors=TRUE
    echo "error: clan ${clanNo}: missing clan data/output folder"
    continue
  }
  [ -f "${clanNo}/ottoapp.json" ] || {
    foundErrors=TRUE
    echo "error: clan ${clanNo}: missing 'ottoapp.json' configuration file'"
    continue
  }
  [ -f "${clanNo}/data/input/ottomap.json" ] || {
    foundErrors=TRUE
    echo "error: clan ${clanNo}: missing 'ottomap.json' configuration file'"
    continue
  }
done
[ -z "${foundErrors}" ] || {
  echo "error: please correct the errors and restart"
  exit 2
}

# parse the turn reports and save the extracts
jq -r 'to_entries | sort_by(.key) | .[] | "\(.key) \(.value)"' clans.json |
while read -r clanNo handle; do
  for turn in 0899-12; do
    inputFile=${clanNo}/data/input/${turn}.${clanNo}.report.docx
    outputFile=${clanNo}/data/input/${turn}.${clanNo}.report.txt
    [ -f "${inputFile}" ] || continue
    [ -f "${outputFile}" ] && {
      [ "${outputFile}" -nt "${inputFile}" ] && {
        echo "${inputFile}: skipping (parsed)"
        continue
      }
      rm -f ${inputFile}
    }
    echo "${inputFile}: parsing..."
    bin/ottoapp run parse report \
      ${inputFile} --output ${outputFile} || {
      echo " error: clan ${clanNo}: handle ${handle}"
      continue
    }
    touch ${outputFile}
    echo "${inputFile}: parsed"
  done
done

exit 0
