#!/usr/bin/env bash

if [ $# -ne 1 ];
then
    echo "Usage: ca_test <ca_id>"
    exit 1
fi

go run . ca ./testData/fakeCA/$1/CA.json