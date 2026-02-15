# OttoApp
A guide to manually load a turn.

## Environment

I have a `.env` file in my working folder to set and export the variables used in this document:

```bash
$ cat .env
export OTTOMAP_REPO=/Users/wraith/JetBrains/tribenet/ottomap
export OTTOAPP_REPO=/Users/wraith/Software/playbymail/ottoapp
export OTTOAPP_ROOT=${OTTOAPP_REPO}/data/production
export OTTOAPP_BIN=${OTTOAPP_ROOT}/bin
export OTTOAPP_GAME=0301
export OTTOAPP_TURN=0900-05
export OTTOAPP_TURNS="0899-12 0900-01 0900-02 0900-03 0900-04 0900-05 0900-06"
export OTTOAPP_CLANS="0500 0506 0507 0508 0509 0511 0513 0514 0515 0518 0520 0538 0542 0544 0555 0556 0557 0582 0594 0633 0668 0669 0711 0720 0729 0866 0912 0913 0925"
export OTTOAPP_FILES=${OTTOAPP_ROOT}/files/${OTTOAPP_GAME}

$ source .env
```

The `OTTOAPP_TURN` variable should be updated as new turns are received from the GM.
I've done that, so you may see files in this document that don't match the turn number above.

# Loading A Turn

1. Wait for the e-mail from the GM.
2. Verify the file names are like 0987.docx.
3. Add new players.
4. Disable dropped players.
5. Download reports. GMail saves them as a single .ZIP file for me.

## Working Area
My working area is a folder in the repository that Git ignores.

```bash
$ cd ${OTTOAPP_ROOT}
$ pwd
/Users/wraith/Software/playbymail/ottoapp/data/production
```

There are three important folders for the input and working files:

```bash
$ ls -l ${OTTOAPP_FILES}
total 0
drwxr-xr-x   31 wraith  staff   992 Dec 14 12:08 ottomap
drwxr-xr-x    2 wraith  staff    64 Dec 16 07:31 report-extracts
drwxr-xr-x  154 wraith  staff  4928 Feb 13 10:10 turn-reports
```

## Creating the Turn Report files
I save turn report files in my working area.

```bash
$ cd ${OTTOAPP_FILES}/turn-reports
$ pwd
/Users/wraith/Software/playbymail/ottoapp/data/production/files/0301/turn-reports
```

I move the download zip file here:

```bash
$ mv ~/Downloads/turnsforottomapfromgame3_1.zip ${OTTOAPP_GAME}.${OTTOAPP_TURN}.zip
ls -l *.tgz *.zip
-rw-r--r--  1 wraith  staff  246647 Dec 15 19:52 0301.0899-12.tgz
-rw-r--r--  1 wraith  staff  362067 Dec 15 19:52 0301.0900-01.tgz
-rw-r--r--  1 wraith  staff  320324 Dec 15 19:52 0301.0900-02.tgz
-rw-r--r--@ 1 wraith  staff  404554 Jan  3 17:30 0301.0900-03.zip
-rw-r--r--@ 1 wraith  staff  341491 Jan 18 07:21 0301.0900-04.zip
```

Note that this breaks things if you forget to update OTTOAPP_TURN.

We need to rename the files from `0987.docx` to `${OTTOAPP_GAME}.{OTTOAPP_TURN}.{clan}.docx`.

```bash
$ ${OTTOAPP_BIN}/rename.sh
```

Return to the working folder and confirm that you have your turn report files with the expected name.


```bash
$ cd ${OTTOAPP_ROOT}
$ ls -l ${OTTOAPP_FILES}/turn-reports/${OTTOAPP_GAME}.${OTTOAPP_TURN}.*.docx
ls -l ${OTTOAPP_FILES}/turn-reports/${OTTOAPP_GAME}.${OTTOAPP_TURN}.*.docx
-rw-r--r--@ 1 wraith  staff  16620 Feb 15 10:18 /Users/wraith/Software/playbymail/ottoapp/data/production/files/0301/turn-reports/0301.0900-06.0500.docx
-rw-r--r--@ 1 wraith  staff  17153 Feb 15 10:18 /Users/wraith/Software/playbymail/ottoapp/data/production/files/0301/turn-reports/0301.0900-06.0508.docx
-rw-r--r--@ 1 wraith  staff  23268 Feb 15 10:18 /Users/wraith/Software/playbymail/ottoapp/data/production/files/0301/turn-reports/0301.0900-06.0509.docx
-rw-r--r--@ 1 wraith  staff  22850 Feb 15 10:18 /Users/wraith/Software/playbymail/ottoapp/data/production/files/0301/turn-reports/0301.0900-06.0511.docx
-rw-r--r--@ 1 wraith  staff  18369 Feb 15 10:18 /Users/wraith/Software/playbymail/ottoapp/data/production/files/0301/turn-reports/0301.0900-06.0518.docx
-rw-r--r--@ 1 wraith  staff  27378 Feb 15 10:18 /Users/wraith/Software/playbymail/ottoapp/data/production/files/0301/turn-reports/0301.0900-06.0520.docx
-rw-r--r--@ 1 wraith  staff  17751 Feb 15 10:18 /Users/wraith/Software/playbymail/ottoapp/data/production/files/0301/turn-reports/0301.0900-06.0538.docx
-rw-r--r--@ 1 wraith  staff  18837 Feb 15 10:18 /Users/wraith/Software/playbymail/ottoapp/data/production/files/0301/turn-reports/0301.0900-06.0542.docx
-rw-r--r--@ 1 wraith  staff  21128 Feb 15 10:18 /Users/wraith/Software/playbymail/ottoapp/data/production/files/0301/turn-reports/0301.0900-06.0544.docx
-rw-r--r--@ 1 wraith  staff  21624 Feb 15 10:18 /Users/wraith/Software/playbymail/ottoapp/data/production/files/0301/turn-reports/0301.0900-06.0555.docx
-rw-r--r--@ 1 wraith  staff  19320 Feb 15 10:18 /Users/wraith/Software/playbymail/ottoapp/data/production/files/0301/turn-reports/0301.0900-06.0556.docx
-rw-r--r--@ 1 wraith  staff  18208 Feb 15 10:18 /Users/wraith/Software/playbymail/ottoapp/data/production/files/0301/turn-reports/0301.0900-06.0582.docx
-rw-r--r--@ 1 wraith  staff  22958 Feb 15 10:18 /Users/wraith/Software/playbymail/ottoapp/data/production/files/0301/turn-reports/0301.0900-06.0633.docx
-rw-r--r--@ 1 wraith  staff  21859 Feb 15 10:18 /Users/wraith/Software/playbymail/ottoapp/data/production/files/0301/turn-reports/0301.0900-06.0668.docx
-rw-r--r--@ 1 wraith  staff  26247 Feb 15 10:18 /Users/wraith/Software/playbymail/ottoapp/data/production/files/0301/turn-reports/0301.0900-06.0720.docx
-rw-r--r--@ 1 wraith  staff  27208 Feb 15 10:18 /Users/wraith/Software/playbymail/ottoapp/data/production/files/0301/turn-reports/0301.0900-06.0912.docx
-rw-r--r--@ 1 wraith  staff  22782 Feb 15 10:18 /Users/wraith/Software/playbymail/ottoapp/data/production/files/0301/turn-reports/0301.0900-06.0925.docx
```

## Extracts

Make sure you're back in the production data folder and run the script.

```bash
$ cd ${OTTOAPP_ROOT}
$  ../../tools/production.sh
```

Errors stop the build. Fix and restart.

A helpful shortcut is
```bash
clan=0500; open files/0301/turn-reports/${OTTOAPP_GAME}.${OTTOAPP_TURN}.${clan}.docx
```

```log
 info: 0301: 0500: running make maps.mk...
extracting data/input/0900-04.0500.report.txt...
make: *** [data/input/0900-04.0500.report.txt] Error 1
error: 0301: 0500: make maps.mk failed
total 8
-rw-r--r--  1 wraith  staff  88 Jan 18 07:47 0900-04.0500.extract.log
```

```bash
$ more files/0301/ottomap/0500/data/errors/0900-04.0500.extract.log
Error: 0301.0900-04.0500.docx:2:27 (88): no match found, expected: "Spring" or "Winter"
```

Success looks like:
```log
import.go:338: sync: import: "0301.0899-12.0925.wxx" 368
import.go:338: sync: import: "0301.0900-01.0925.wxx" 369
import.go:338: sync: import: "0301.0900-02.0925.wxx" 370
import.go:338: sync: import: "0301.0900-03.0925.wxx" 371
import.go:338: sync: import: "0301.0900-04.0925.wxx" 372
sync.go:219: sync: import: map-files: completed in 47.8295ms
 info: synced map files to the database
```

Verify that the database is updated and doesn't need to be compacted:

```bash
$ ll db
total 8224 -rw-r--r--  1 wraith  staff   3.9M Jan 18 08:21 ottoapp.db
```

## Test a migration to stage

```bash
$ ssh ottomap

root# systemctl stop ottoapp-stage
```

```bash
 production git:(main) ✗ scp db/ottoapp.db ottomap:/var/www/stage/ottoapp/data/
ottoapp.db      100% 3972KB   1.0MB/s   00:03
➜  production git:(main) ✗
```

Restart the stage server
```bash
$ ssh ottomap systemctl restart ottoapp-stage
```

Login to https://ottomap-stage.playbymailgames.com/ and verify a map.

Then, if that works, push the file to production.

```bash
root# systemctl stop ottoapp
root# cd /var/www/prod/ottoapp/data/
root# ls -l
-rw-r--r-- 1 ottopb ottopb 4067328 Jan 18 13:32 ottoapp.db
-rw-r--r-- 1 ottopb ottopb   32768 Jan 18 13:36 ottoapp.db-shm
-rw-r--r-- 1 ottopb ottopb   32992 Jan 18 13:36 ottoapp.db-wal
root# rm ottoapp.db*
root# cp -p /var/www/stage/ottoapp/data/ottoapp.db /var/www/prod/ottoapp/data/
root# systemctl start ottoapp  
```

And login to https://ottomap.playbymailgames.com/ and verify a map.
