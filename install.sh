#!/bin/bash

# fail on any error
set -e

# install variables
export APP_NAME=cgen
export BUCKET_URL=https://s3-us-west-2.amazonaws.com/github.techdecaf.io
export INSTALL_PATH=${INSTALL_PATH:-"/usr/local/bin"}
export INSTALL_VERSION=${INSTALL_VERSION:-"latest"}

while getopts v:p: flag
do
  case "${flag}" in
    v) INSTALL_VERSION=${OPTARG};;
    p) INSTALL_PATH=${OPTARG};;
  esac
done

# export os as lowercase string
export OS=$(echo $(uname) | tr '[:upper:]' '[:lower:]')

if [[ "$OS" == "" ]]; then
    echo ERROR: failed to determine operating system version
    exit 1
fi

export APP_BINARY=/tmp/$APP_NAME
export LATEST_STABLE=$BUCKET_URL/$APP_NAME/$INSTALL_VERSION/$OS/$APP_NAME

echo '[Installed] '$APP_NAME version: $($APP_NAME --version)
echo '[Downloading]' $LATEST_STABLE && curl -fsSLo $APP_BINARY $LATEST_STABLE
echo '[Installing]' $APP_NAME && chmod +x $APP_BINARY && mv $APP_BINARY $INSTALL_PATH
echo '[Validation]' $APP_NAME version: $($APP_NAME --version)