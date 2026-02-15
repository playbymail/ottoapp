#!/bin/bash
############################################################################
# runall.sh - synchronize the database and file system, then build maps
############################################################################
# set up and export environment for child processes
export OTTOAPP_ROOT=/Users/wraith/Software/playbymail/ottoapp
cd "${OTTOAPP_ROOT}" || {
  echo "error: unable to set def to OTTAPP_ROOT"
  exit 2
}
export OTTOMAP_ROOT=/Users/wraith/Jetbrains/tribenet/ottomap
cd "${OTTOMAP_ROOT}" || {
  echo "error: unable to set def to OTTMAP_ROOT"
  exit 2
}
export TN_ROOT=/Users/wraith/Software/playbymail/ottoapp/data/production
cd "${TN_ROOT}" || {
  echo "error: unable to set def to TN_ROOT"
  exit 2
}
############################################################################
# local environment
dbPath=./db
ottoapp=${TN_ROOT}/bin/ottoapp
ottomap=${TN_ROOT}/bin/ottomap
############################################################################
# command line flags. these are listed in the order they execute.
showHelp=NO
quiet=NO
verbose=NO
debug=NO
rebuildOttoapp=YES
rebuildOttomap=YES
ottoappLogging=
rebuildDatabase=YES
importConfig=YES
importUsers=YES
importGames=YES
rebuildClanFolders=YES
rebuildMakefiles=YES
renderMaps=YES
syncReportFiles=YES
syncReportExtractFiles=YES
syncMapKeyFiles=YES
syncMapFiles=YES

############################################################################
# process command line flags
for opt in "$@"; do
  shift
  case "${opt}" in
  -h | --help                    ) showHelp=YES ;;
  -c | --import-config           ) importConfig=YES ;;
  -d | --debug                   ) debug=YES;;
  -g | --import-games            ) importGames=YES ;;
  -q | --quiet                   ) quiet=YES; verbose=NO ;;
  -u | --import-users            ) importUsers=YES ;;
  -v | --verbose                 ) quiet=NO; verbose=YES ;;
  -a | --rebuild-ottoapp         ) rebuildOttoapp=YES;;
  -A | --no-rebuild-ottoapp      ) rebuildOttoapp=NO ;;
  -f | --rebuild-clan-folders    ) rebuildClanFolders=YES ;;
  -F | --no-rebuild-clan-folders ) rebuildClanFolders=NO ;;
  -k | --rebuild-makefiles       ) rebuildMakefiles=YES;;
  -K | --no-rebuild-makefiles    ) rebuildMakefiles=NO ;;
  -m | --rebuild-ottomap         ) rebuildOttomap=YES;;
  -M | --no-rebuild-ottomap      ) rebuildOttomap=NO ;;
  --rebuild-database)
    rebuildDatabase=YES
    importConfig=YES
    importUsers=YES
    importGames=YES ;;
  --render-map-files)
    renderMaps=YES ;;
  --sync-report-files)
    syncReportFiles=YES ;;
  --sync-report-extract-files)
    syncReportExtractFiles=YES ;;
  --sync-map-key-files)
    syncMapKeyFiles=YES ;;
  --sync-map-files)
    syncMapFiles=YES ;;
  *)
    echo "error: unknown option '${opt}'"
    exit 1;;
  esac
done

echo " info: dbPath                 == '${dbPath}'"
echo "     : quiet                  == '${quiet}'"
echo "     : verbose                == '${verbose}'"
echo "     : debug                  == '${debug}'"
echo "     : rebuildOttoapp         == '${rebuildOttoapp}'"
echo "     : rebuildOttomap         == '${rebuildOttomap}'"
echo "     : rebuildDatabase        == '${rebuildDatabase}'"
echo "     : importConfig           == '${importConfig}'"
echo "     : importUsers            == '${importUsers}'"
echo "     : importGames            == '${importGames}'"
echo "     : rebuildClanFolders     == '${rebuildClanFolders}'"
echo "     : rebuildMakefiles       == '${rebuildMakefiles}'"
echo "     : renderMaps             == '${renderMaps}'"
echo "     : syncReportFiles        == '${syncReportFiles}'"
echo "     : syncReportExtractFiles == '${syncReportExtractFiles}'"
echo "     : syncMapKeyFiles        == '${syncMapKeyFiles}'"
echo "     : syncMapFiles           == '${syncMapFiles}'"

[ "${showHelp}" == "YES" ] && {
  echo "usage: ./production.sh [options]"
  echo "  opt: --help                      | -h :  show this text"
  echo "       --rebuild-database               :  rebuild the database"
  echo "       --import-config             | -c :  import the configuration file"
  echo "       --import-games              | -g :  import the configuration file"
  echo "       --import-users              | -u :  import the users file"
  echo "       --sync-report-files              :  import turn report files"
  echo "       --sync-report-extract-files      :  import report extract files"
  echo "       --sync-map-key-files             :  import map keys"
  echo "       --sync-map-files                 :  import map files"
  echo "       --rebuild-ottoapp           | -a :  rebuild ottoapp before running"
  echo "       --no-rebuild-ottoapp        | -A :  do not rebuild ottoapp"
  echo "       --rebuild-ottomap           | -m :  rebuild ottomap before running"
  echo "       --no-rebuild-ottomap        | -M :  do not rebuild ottomap"
  echo "       --debug                     | -d :  enable debug flag in ottoapp/ottomap"
  echo "       --quiet                     | -q :  enable quiet flag in ottoapp/ottomap"
  echo "       --verbose                   | -v :  enable verbose flag in ottoapp/ottomap"
  exit 1
}

# convert quiet, verbose, and debug flags for ottoapp
if [ "${quiet}" == "YES" ]; then
  ottoappLogging="--quiet"
elif [ "${verbose}" == "YES" ]; then
  ottoappLogging="--verbose"
fi
if [ "${debug}" == "YES" ]; then
  ottoappLogging="--debug ${ottoappLogging}"
fi

############################################################################
# check environment
[ -n "${OTTOAPP_ROOT}" ] || {
  echo "error: OTTAPP_ROOT is not set and exported"
  exit 2
}
[ -n "${OTTOMAP_ROOT}" ] || {
  echo "error: OTTOMAP_ROOT is not set and exported"
  exit 2
}
[ -n "${TN_ROOT}" ] || {
  echo "error: TN_ROOT is not set and exported"
  exit 2
}
cd "${TN_ROOT}" || {
  echo "error: unable to set def to TN_ROOT"
  exit 2
}
[ -d "${dbPath}" ] || {
  echo " info: TN_ROOT             == '${TN_ROOT}'"
  echo "     : dbPath              == '${dbPath}'"
  echo "error: dbPath is not a directory"
  exit 2
}

############################################################################
# build ottoapp
if [ "${rebuildOttoapp}" == "YES" ]; then
  echo " info: rebuilding '${ottoapp}'..."
  echo " info: OTTOAPP_ROOT        == '${OTTOAPP_ROOT}'"
  cd ${OTTOAPP_ROOT} || exit 2
  go build -o ${ottoapp} ./cmd/ottoapp || exit 2
  echo " info: rebuilt '${ottoapp}'"
else
  echo " skip: rebuilding '${ottoapp}'"
fi
echo " info: OTTOAPP_VERSION     == '$( ${ottoapp} version )'"

############################################################################
# build ottomap
if [ "${rebuildOttomap}" == "YES" ]; then
  echo " info: rebuilding '${ottomap}'..."
  echo " info: OTTOMAP_ROOT        == '${OTTOMAP_ROOT}'"
  cd ${OTTOMAP_ROOT} || exit 2
  go build -o ${ottomap} || exit 2
  echo " info: rebuilt ${ottomap}"
  echo " info: rebuilding '${ottomap}-linux'..."
  echo " info: OTTOMAP_ROOT        == '${OTTOMAP_ROOT}'"
else
  echo " skip: rebuilding '${ottomap}'"
fi
echo " info: OTTOMAP_VERSION     == '$( ${ottomap} version )'"

############################################################################
# rebuild the database
if [ "${rebuildDatabase}" == "YES" ]; then
  echo " info: rebuilding database '${dbPath}'..."
  # force the user to set a lot a flags if we're rebuilding the database
  [ "${importConfig}" == "YES" ] || {
    echo "error: you must set --load-config when rebuilding the database"
  }
  [ "${importUsers}" == "YES" ] || {
    echo "error: you must set --load-users when rebuilding the database"
  }
  cd "${TN_ROOT}" || {
    echo "error: unable to set def to TN_ROOT"
    exit 2
  }
  ${ottoapp} ${ottoappLogging} db init --overwrite "${dbPath}" || {
    echo "error: error rebuilding database"
    exit 2
  }
  ${ottoapp} --db "${dbPath}" db version
  echo " info: rebuilt database '${dbPath}'"
else
  echo " skip: rebuilding database '${dbPath}'"
fi

############################################################################
# import our configuration JSON file
if [ "${importConfig}" == "YES" ]; then
  echo " info: importing ottoapp configuration './config/ottoapp.json'..."
  cd "${TN_ROOT}" || {
    echo "error: unable to set def to TN_ROOT"
    exit 2
  }
  ${ottoapp} ${ottoappLogging} --db "${dbPath}" sync import ottoapp-config-file ./config/ottoapp.json || {
    echo "error: import ottoapp-config-file failed"
    exit 2
  }
  echo " info: imported ottoapp configuration './config/ottoapp.json'"
else
  echo " skip: importing ottoapp configuration './config/ottoapp.json'"
fi

############################################################################
# import our user data JSON file
if [ "${importUsers}" == "YES" ]; then
  echo " info: importing users './config/users.json'..."
  cd "${TN_ROOT}" || {
    echo "error: unable to set def to TN_ROOT"
    exit 2
  }
  ${ottoapp} ${ottoappLogging} --db "${dbPath}" sync import users ./config/users.json || {
    echo "error: import users failed"
    exit 2
  }
  echo " info: imported users './config/users.json'"
else
  echo " skip: importing users './config/users.json'"
fi

############################################################################
# import our game data JSON file
if [ "${importGames}" == "YES" ]; then
  echo " info: importing games './config/games.json'..."
  cd "${TN_ROOT}" || {
    echo "error: unable to set def to TN_ROOT"
    exit 2
  }
  ${ottoapp} ${ottoappLogging} --db "${dbPath}" sync import games ./config/games.json || {
    echo "error: import games failed"
    exit 2
  }
  echo " info: imported games './config/games.json'"
else
  echo " skip: importing games './config/games.json'"
fi
############################################################################
# create any missing clan directories in the ottomap folder
if [ "${rebuildClanFolders}" == "YES" ]; then
  echo " info: building clan directories for ottomap..."
  cd "${TN_ROOT}" || {
    echo "error: unable to set def to TN_ROOT"
    exit 2
  }
  for game in 0300 0301; do
    for file in files/${game}/turn-reports/${game}.????-??.0???.docx; do
      [ -f "${file}" ] || continue
      # extract the clan from the file name
      fileSansSuffix=${file%.docx}
      clanNo=${fileSansSuffix##*.}
      # make the data directories (use -p so we don't error on rebuild)
      mkdir -p files/${game}/ottomap/${clanNo}/data/{errors,input,logs,output}
    done
  done
  echo " info: built clan directories for ottomap"
else
  echo " skip: building clan directories for ottomap"
fi

############################################################################
# rebuild the makefile
if [ "${rebuildMakefiles}" == "YES" ]; then
  echo " info: rebuilding makefiles..."
  cd "${TN_ROOT}" || {
    echo "error: unable to set def to TN_ROOT"
    exit 2
  }
  ${ottoapp} ${ottoappLogging} --db "${dbPath}" generate makefile files/0301 || {
    echo "error: generate makefile failed"
    exit 2
  }
  echo " info: rebuilt makefiles"
else
  echo " skip: rebuilding makefiles"
fi

############################################################################
# run the makefiles to render the maps
if [ "${renderMaps}" == "YES" ]; then
  echo " info: rendering map files..."
  cd "${TN_ROOT}" || {
    echo "error: unable to set def to TN_ROOT"
    exit 2
  }
  for game in 0300 0301; do
    [ -d "${TN_ROOT}/files/${game}/ottomap" ] || continue
    cd ${TN_ROOT}/files/${game}/ottomap || {
      echo "error: unable to set def to TN_ROOT/files/${game}/ottomap"
      exit 2
    }
    for clan in 0???; do
      [ -f "${TN_ROOT}/files/${game}/ottomap/${clan}/maps.mk" ] || continue
      cd ${TN_ROOT}/files/${game}/ottomap/${clan} || exit 2
      echo " info: ${game}: ${clan}: running make maps.mk..."
      make -f maps.mk || {
        echo "error: ${game}: ${clan}: make maps.mk failed"
        ls -l data/errors
        exit 2
      }
    done
  done
  echo " info: rendered map files"
else
  echo " skip: rendering map files"
fi

############################################################################
# sync report files to the database
if [ "${syncReportFiles}" == "YES" ]; then
  echo " info: syncing report files to the database..."
  cd "${TN_ROOT}" || {
    echo "error: unable to set def to TN_ROOT"
    exit 2
  }
  ${ottoapp} ${ottoappLogging} --db "${dbPath}" sync import turn-report-files ./ --verbose|| {
    echo "error: sync import turn-report-files failed"
    exit 2
  }
  echo " info: synced report files to the database"
else
  echo " skip: syncing report files to the database"
fi

############################################################################
# sync report extract files to the database
if [ "${syncReportExtractFiles}" == "YES" ]; then
  echo " info: syncing report extract files to the database..."
  cd "${TN_ROOT}" || {
    echo "error: unable to set def to TN_ROOT"
    exit 2
  }
  ${ottoapp} ${ottoappLogging} --db "${dbPath}" sync import report-extract-files ./ --verbose || {
    echo "error: sync import report-extract-files failed"
    exit 2
  }
  echo " info: syncing report extract files to the database"
else
  echo " skip: syncing report extract files to the database"
fi

#############################################################################
## sync map key files
if [ "${syncMapKeyFiles}" == "YES" ]; then
  echo " info: syncing map key files..."
  #./upload-map-keys.sh                  || exit 2
  echo " info: synced map key files"
else
  echo " skip: syncing map key files"
fi

############################################################################
# sync map files to the database
if [ "${syncMapFiles}" == "YES" ]; then
  echo " info: syncing map files to the database..."
  cd "${TN_ROOT}" || {
    echo "error: unable to set def to TN_ROOT"
    exit 2
  }
  ${ottoapp} ${ottoappLogging} --db "${dbPath}" sync import map-files ./ --verbose || {
    echo "error: sync import map-files failed"
    exit 2
  }
  echo " info: synced map files to the database"
else
  echo " skip: syncing map files to the database"
fi

exit 0
