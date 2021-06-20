#!/bin/bash

for shard in localhost:8080 localhost:8081; do 
    for i in {1..1000}; do 
        echo "http://$shard/set?key=key-$i&value=value-$i" &
        curl "http://$shard/set?key=key-$i&value=value-$i"
    done
done
