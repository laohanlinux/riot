# riot

<center>![](doc/riot.jpg)</center>

Riot is a distributed key/value system basing raft algorithm, leveldb、 boltdb and bitCask(in fucture) backend store!!!

## Install Doc

- install grpc

```sh
go get -u google.golang.org/grpc
go get -a github.com/golang/protobuf/protoc-gen-go
```

*Notice*:

> if can not compile it, please do that (may sure protoc is installed): 
> cd rpc/pb && protoc --go_out=plugins=grpc:. op.proto

- build riot

```css
go build riot
```

- start a cluster

```shell
cd tool
bash cluster1.sh
```

## Api 

in doc directory

## about design detail 

[link](https://laohanlinux.github.io/2016/04/25/%E4%BD%BF%E7%94%A8raft%E7%AE%97%E6%B3%95%E5%BF%AB%E7%86%9F%E6%9E%84%E5%BB%BA%E4%B8%80%E4%B8%AA%E5%88%86%E5%B8%83%E5%BC%8F%E7%9A%84key-value%E7%B3%BB%E7%BB%9F/)
