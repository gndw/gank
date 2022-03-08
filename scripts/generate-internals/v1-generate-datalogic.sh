#!/bin/sh

PACKAGE_FOLDER=$1
PACKAGE_NAME=$2

cd "$(dirname "$0")"

sh generate.sh v1/internals/datalogics/$PACKAGE_FOLDER $PACKAGE_NAME Datalogic d datalogic