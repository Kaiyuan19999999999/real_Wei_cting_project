#!/usr/bin/env bash

if [ $# -ne 1 ];
then
    echo "Usage: gossiper_test <gossiper_id>"
    exit 1
fi

go run . gossiper ./testData/gossiperNetworkTest/gossiper_pub_config.json ./testData/gossiperNetworkTest/$1/gossiper_priv_config.json ./testData/gossiperNetworkTest/$1/gossiperCrypto.json $1