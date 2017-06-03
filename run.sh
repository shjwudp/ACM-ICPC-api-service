#!/bin/env bash

# access workspace
SCRIPT_DIR=$(
	cd $(dirname $0)
	pwd
)
cd ${SCRIPT_DIR}

go run main.go --config ./conf.json