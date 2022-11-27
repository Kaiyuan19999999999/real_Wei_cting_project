#!/usr/bin/env bash

if [ $# -ne 1 ];
then
    echo "Usage: logger_test <logger_id>"
    exit 1
fi

go run . logger ./testData/fakeLogger/$1/logger.json