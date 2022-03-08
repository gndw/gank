#!/bin/sh

ROOT_RELATIVE_FOLDER=$1
PACKAGE_NAME=$2

cd "$(dirname "$0")"

sh generate.sh $ROOT_RELATIVE_FOLDER $PACKAGE_NAME Repository r repository