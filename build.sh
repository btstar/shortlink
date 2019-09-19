#!/bin/bash

WORKDIR=$(pwd)

cd $WORKDIR && go build -tags netgo -v .
