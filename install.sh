#!/bin/bash

# export os as lowercase string
export OS=$(echo $(uname) | tr '[:upper:]' '[:lower:]')

if [[ "$OS" == "" ]]; then
    echo ERROR: failed to determine operating system version
    exit 1
fi

# CONSTS
export APP_NAME=cgen
export BUCKET_URL=http://github.techdecaf.io
export INSTALL_PATH=/usr/local/bin

export APP_BINARY=/tmp/$APP_NAME
export LATEST_STABLE=$BUCKET_URL/$APP_NAME/latest/$OS/$APP_NAME

echo '[Installed] '$APP_NAME version: $($APP_NAME --version)
echo '[Downloading]' $LATEST_STABLE && curl -o $APP_BINARY $LATEST_STABLE
echo '[Installing]' $APP_NAME && chmod +x $APP_BINARY && mv $APP_BINARY $INSTALL_PATH
echo '[Validation]' $APP_NAME version: $($APP_NAME --version)