#!/bin/bash

# Print usage string
function showUsage {
  echo "usage: $(basename $0) [-f] [-b] [-o output_folder] [-c]"
  echo -e "\t-f front-end force: force building the front-end"
  echo -e "\t-b back-end force: force building the back-end"
  echo -e "\t-o output_folder: location of final build project"
  echo -e "\t-c clean: remove all files from the output folder prior to building"
}

# Get absolute path
function fullPath {
  TARGET_FILE=$1

  cd `dirname $TARGET_FILE`
  TARGET_FILE=`basename $TARGET_FILE`

  PHYS_DIR=`pwd -P`
  RESULT=$PHYS_DIR/$TARGET_FILE
  RESULT=${RESULT%/.}
  echo $RESULT
}

# Deep search for the latest modification
function getLatestDate {
  LATEST=0
  while read -r source
  do
    CURRENT=$(date -r $source +%s)
    if [ $CURRENT -gt $LATEST ]
    then
      LATEST=$CURRENT
    fi
  done <<< "$1"

  echo $LATEST
}

# Parse arguments
while getopts "fbo:c" opt
do
  case "$opt" in
    f)
      BUILD_FRONT=true
      ;;
    b)
      BUILD_BACK=true
      ;;
    o)
      OUTPUT_FOLDER="$OPTARG"
      ;;
    c)
      CLEAN=true
      ;;
    *)
      echo "[31mInvalid flag[m"
      showUsage
      exit 1
      ;;
  esac
done

# Check GO environment
if [ ! "$GOPATH" ]
then
  echo "[31mCould not find GOPATH[m"
  echo "[31mTry issuing:[m"
  echo "export GOPATH=<Path to where GO projects are>"
  exit 1
fi
if [[ ! $(go version) ]]
then
  echo "[31mCould not find GO[m"
  exit 1
fi

BASE_DIR=$(dirname $(fullPath $0))
PROJECT=${BASE_DIR##*/}
echo "Building $PROJECT"

# Check output folder
if [ "$OUTPUT_FOLDER" ]
then
  OUTPUT_FOLDER=$(fullPath "$OUTPUT_FOLDER")
  if [ -f "$OUTPUT_FOLDER" ]
  then
    echo "[31mThe given output folder is invalid: $OUTPUT_FOLDER[m"
    showUsage
    exit 1
  fi
else
  OUTPUT_FOLDER="$BASE_DIR/build"
fi

echo "Building at $OUTPUT_FOLDER"
pushd "$BASE_DIR" &> /dev/null

# Clean, if required
if [ $CLEAN ]
then
  rm -r "$OUTPUT_FOLDER"
fi

# Check if should build front-end
if [[ ! "$BUILD_FRONT" ]]
then
  if [ -d "$OUTPUT_FOLDER/web" ]
  then
    SOURCES=$(find $(git ls-files web/) -type f -not -name '.gitignore')
    SOURCE_DATE=$(getLatestDate "$SOURCES")
    BUILD_DATE=$(date -r "$OUTPUT_FOLDER/web" +%s)

    if [ $SOURCE_DATE -gt $BUILD_DATE ]
    then
      BUILD_FRONT=true
    fi
  else
    BUILD_FRONT=true
  fi
fi

# Build front-end
if [[ "$BUILD_FRONT" ]]
then
  echo "[32mBuilding front-end..[m"
  pushd web &> /dev/null
  npm install &&  npm run build

  if [ $? -eq 0 ]
  then
    echo "[32mDone[m"
    popd &> /dev/null
  else
    echo "[31mFailed[m"
    popd &> /dev/null
    popd &> /dev/null
    exit 1
  fi
else
  echo "[32mSkipping front-end..[m"
fi

# Check if should build back-end
if [[ ! "$BUILD_BACK" ]]
then
  if [ -f "$OUTPUT_FOLDER/$PROJECT" ]
  then
    SOURCES=$(find $(git ls-files) -type f -name '*.go')
    SOURCE_DATE=$(getLatestDate "$SOURCES")
    BUILD_DATE=$(date -r "$OUTPUT_FOLDER/$PROJECT" +%s)

    if [ $SOURCE_DATE -gt $BUILD_DATE ]
    then
      BUILD_BACK=true
    fi
  else
    BUILD_BACK=true
  fi
fi

# Build back-end
if [[ "$BUILD_BACK" ]]
then
  echo "[32mBuilding back-end[m"
  go get && go install

  if [ $? -eq 0 ]
  then
    echo "[32mDone[m"
  else
    echo "[31mFailed[m"
    popd &> /dev/null
    exit 1
  fi
else
  echo "[32mSkipping back-end..[m"
fi

# Copy files
echo "[32mCopying files..[m"

if [ ! -d "$OUTPUT_FOLDER" ]
then
  mkdir "$OUTPUT_FOLDER"
  cp "$GOPATH/bin/$PROJECT" "$OUTPUT_FOLDER"/.
  cp -r web/build "$OUTPUT_FOLDER/web"
  cp secrets/* "$OUTPUT_FOLDER/." 2> /dev/null
  cp *.conf "$OUTPUT_FOLDER/." 2> /dev/null

else
  cp secrets/ "$OUTPUT_FOLDER/." 2> /dev/null
  cp *.conf "$OUTPUT_FOLDER/." 2> /dev/null

  if [ $BUILD_BACK ]
  then
    cp "$GOPATH/bin/$PROJECT" "$OUTPUT_FOLDER"/.
  fi

  if [ $BUILD_FRONT ] || [ ! -d "$OUTPUT_FOLDER/web" ]
  then
    cp -r web/build "$OUTPUT_FOLDER/web"
  fi
fi

popd &> /dev/null
