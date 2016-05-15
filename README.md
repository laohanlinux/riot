# riot

<center>![](doc/riot.jpg)</center>

Riot is a distributed key/value system basing raft algorithm, leveldb and bitCask backend storage!!!


## Install Doc

- install grpc

```
go get -u google.golang.org/grpc
go get -a github.com/golang/protobuf/protoc-gen-go
```

*Notice*:

> if can not compile it, please do that (may sure protoc is installed): 
> cd rpc/pb && protoc --go_out=plugins=grpc:. op.proto

- build riot

```
go build riot
```

- start a cluster

```
cd tool
bash cluster1.sh 1
```
