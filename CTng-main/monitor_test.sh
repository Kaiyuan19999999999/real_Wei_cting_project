#!/usr/bin/env bash

if [ $# -ne 1 ];
then
    echo "Usage: monitor_test <monitor_id>"
    exit 1
fi

go run . monitor ./testData/monitorNetworkTest/monitor_pub_config.json ./testData/monitorNetworkTest/$1/monitor_priv_config.json ./testData/monitorNetworkTest/$1/monitorCrypto.json $1