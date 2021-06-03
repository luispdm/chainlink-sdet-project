#!/bin/bash

VER="v0.6.6/"
names=()
cd ./contracts/ethereum/$VER/abi
for FILE in *; do
    names+=($(echo "${FILE%%.*}"))
done

cd ../../
for elem in "${names[@]}"; do
    if [ ! -f "$elem.go" ]; then
        abigen --bin=$VER/bin/$elem.bin --abi=$VER/abi/$elem.abi --pkg=$elem --out=$elem.go
        sed -i "s/package $elem/package ethereum/" $elem.go
    else
        echo "'$elem.go' already exists, skipping generation of file"
    fi
done
