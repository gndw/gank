#!/bin/sh

# example call
# sh generate.sh v1/internals/datalogics/pkg pkg Datalogic d datalogic

#define parameters which are passed in.
# v1/internals/datalogic/dummy
ROOT_RELATIVE_FOLDER=$1
# auth
PACKAGE_NAME=$2
# Datalogic
TYPE=$3
# d
MINI_TYPE=$4
# datalogic
MINI_TYPE_NAME=$5

FINALPATH="../../"$ROOT_RELATIVE_FOLDER

# create parent folder
mkdir -p $FINALPATH
# create interface
sed -e "s;%PACKAGE_NAME%;$PACKAGE_NAME;g" -e "s;%TYPE%;$TYPE;g" -e "s;%MINI_TYPE%;$MINI_TYPE;g" interface.template.txt > $FINALPATH/$MINI_TYPE_NAME.go
# create impl folder
mkdir -p $FINALPATH/impl
# create init file
sed -e "s;%PACKAGE_NAME%;$PACKAGE_NAME;g" -e "s;%TYPE%;$TYPE;g" -e "s;%MINI_TYPE%;$MINI_TYPE;g" -e "s;%ROOT_RELATIVE_FOLDER%;$ROOT_RELATIVE_FOLDER;g" init.template.txt > $FINALPATH/impl/init.go
# create func file
sed -e "s;%PACKAGE_NAME%;$PACKAGE_NAME;g" -e "s;%TYPE%;$TYPE;g" -e "s;%MINI_TYPE%;$MINI_TYPE;g" func.template.txt > $FINALPATH/impl/func.go