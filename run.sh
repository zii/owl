#!/bin/bash

set -e

# run one/many
for arg in "$@"
do
	if [[ $arg == 1 ]]; then
	  echo "run.."
		go build -o ./owl
    ./owl
  elif [[ $arg == 2 ]]; then
    echo "build windows.."
    CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o ./owl.exe
	else
		echo unknown argument: $arg
	fi
done