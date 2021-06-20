#!/bin/bash

set -e

trap "killall dist-kv" SIGINT

cd $(dirname $0)

killall dist-kv || true
sleep 1

go install -v

dist-kv --db-location=$PWD/north.db --config-file=$PWD/sharding.toml --shard=north --http-addr=localhost:8080 &
dist-kv --db-location=$PWD/east.db --config-file=$PWD/sharding.toml --shard=east --http-addr=localhost:8081 &
dist-kv --db-location=$PWD/west.db --config-file=$PWD/sharding.toml --shard=west --http-addr=localhost:8082 &
dist-kv --db-location=$PWD/south.db --config-file=$PWD/sharding.toml --shard=south --http-addr=localhost:8083 &

wait
