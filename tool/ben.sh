#!/bin/bash

for((i=1;i<=1024;i++));do
    url="http://localhost:$1/riot/key/$i"
    echo $url
    curl -XPOST "$url" -d $i
done 


for((i=1;i<=1024;i++));do
    url="http://localhost:$1/riot/key/$i?qs=$2"
    echo $url
    curl -v "$url"
done
