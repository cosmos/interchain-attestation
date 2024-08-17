#!/bin/bash
set -eE -o functrace

failure() {
  local lineno=$1
  local msg=$2
  echo "Failed at $lineno: $msg"
}
trap 'failure ${LINENO} "$BASH_COMMAND"' ERR

source scripts/common-serve-env.sh

ROOT_DIR=/tmp
MOCK_DA_DIR=$ROOT_DIR/mock-da
LOG_FILE_PATH=$MOCK_DA_DIR/mock-da.log

echo "Stopping and removing any existing mock-da container and previous data..."

# Check if container exists
if [ "$(docker ps -q -f name=mockda)" ]; then
  if [ "$(docker inspect -f '{{.State.Running}}' mockda)" == "true" ]; then
    docker container stop mockda 2>/dev/null
  fi

  docker container rm mockda 2>/dev/null
fi



if [ -d $MOCK_DA_DIR ]; then
  echo "Removing previous data..."
  rm -rf $MOCK_DA_DIR &> /dev/null
fi

# Add directories for chain(s), exit if an error occurs
if ! mkdir -p $MOCK_DA_DIR 2>/dev/null; then
    echo "Failed to create chain folder. Aborting..."
    exit 1
fi

# Start a new mockda container
docker run -p $MOCK_DA_PORT:$MOCK_DA_PORT --name mockda mock-da:local /usr/bin/mock-da -listen-all > $LOG_FILE_PATH 2>&1 &

sleep 1

# Check if the container is running
if [ "$(docker inspect -f '{{.State.Running}}' mockda)" == "true" ]; then
  echo "Mock-da container is running"
else
  echo "Mock-da container is not running"
  exit 1
fi
