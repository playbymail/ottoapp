#!/bin/bash
###########################################################################
# rebuild the alpha database
###########################################################################
OTTOREPO=/Users/wraith/Software/playbymail/ottoapp
OTTOPATH=/Users/wraith/Software/playbymail/ottoapp/data/alpha
OTTOBIN="${OTTOPATH}/bin"
ottoapp="${OTTOBIN}/ottoapp"
ottodb="${OTTOPATH}/db"

###########################################################################
# verify paths
cd "${OTTOPATH}" || {
  echo "error: can not set def to OTTOPATH"
  exit 2
}

###########################################################################
# build the executable
cd "${OTTOREPO}" || {
  echo "error: can not set def to OTTOREPO"
  exit 2
}
go build -o "${ottoapp}" ./cmd/ottoapp || {
  echo "error: build failed for ottoapp"
  exit 2
}
"${ottoapp}" version || {
  echo "error: ottoapp version failed"
  exit 2
}

###########################################################################
#
cd "${OTTOPATH}" || {
  echo "error: can not set def to OTTOPATH"
  exit 2
}
"${ottoapp}" db init --overwrite "${ottodb}" || {
  echo "error: db init '${ottodb}' failed"
  exit 2
}

"${ottoapp}" db --db "${ottodb}" version

###########################################################################
#
exit 0
