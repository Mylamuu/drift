#!/bin/sh

rm -rf bin 2> /dev/null
mkdir bin
go build -o bin/drift ./cmd

