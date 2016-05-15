#!/bin/bash

for((i=1;i<=1024;i++));do
    url="http://localhost:8082/riot/key/$i"
    echo $url
    curl -XPOST "$url" -d $i
done 


for((i=1;i<=1024;i++));do
    url="http://localhost:8081/riot/key/$i?qs=$1"
    echo $url
    curl -v "$url"
done
