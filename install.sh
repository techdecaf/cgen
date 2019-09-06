#!/bin/bash


# export os as lowercase string
export OS=$(echo $(uname) | tr '[:upper:]' '[:lower:]')
export APP_NAME=cgen
export APP_VERSION=$1

if [[ "$OS" == "" ]]; then
    echo ERROR: failed to determine operating system version
    exit 1
fi

if [[ "$APP_VERSION" == "" ]]; then
    export APP_VERSION="latest"
fi

# CONSTS
export BUCKET_URL=http://github.techdecaf.io
export INSTALL_PATH=/usr/local/bin

export APP_BINARY=/tmp/$APP_NAME
export APPLICATION=$BUCKET_URL/$APP_NAME/$APP_VERSION/$OS/$APP_NAME

echo '[Installed] '$APP_NAME version: $($APP_NAME --version)
echo '[Downloading]' $APPLICATION #&& curl -o $APP_BINARY $APPLICATION
echo '[Installing]' $APP_NAME #&& chmod +x $APP_BINARY && mv $APP_BINARY $INSTALL_PATH
echo '[Validation]' $APP_NAME version: $($APP_NAME --version)