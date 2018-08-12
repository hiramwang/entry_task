#!/usr/bin/env bash

if [ "$1" == "clean" ]; then                                                                                                                                                                                                                                            
    rm -rf bin/* pkg/* 
    exit
fi

CURDIR=`pwd`
export GOPATH=$GOPATH":"$CURDIR
echo $GOPATH

echo "building entry_task"                                                                                                                                                                                                                                                     
cd src
go build -o entry_task
cd -
mkdir -p bin
mv src/entry_task bin/
