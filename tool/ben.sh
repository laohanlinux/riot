#!/bin/bash

for((i=1;i<=10024;i++));do 
    url="http://localhost:8080/riot?key=$i"
    echo $url
    curl -XPOST "$url" -d $i
done 


for((i=1;i<=10024;i++));do 
    url="http://localhost:8080/riot?key=$i"
    echo $url
    curl -v "$url"
done