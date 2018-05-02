#!/bin/sh

cd $(dirname $0)

BASE="https://net-mozaws-prod-delivery-inventory-us-east-1.s3.amazonaws.com/public/inventories/net-mozaws-prod-delivery-archive/delivery-archive/data"

if [ ! -e "./testdata" ]; then
    mkdir ./testdata
fi

while read -r filename size; do
    if [ -e "./testdata/$filename" ]; then
        if [ $(wc -c "./testdata/$filename" | awk '{print $1}') -ne $size ]; then
            echo "Wrong size ... download again: $filename ($size)"
            curl -o "./testdata/$filename" "$BASE/$filename"
        fi
    else
        curl -o "./testdata/$filename" "$BASE/$filename"
    fi
done < ./1gb.txt
