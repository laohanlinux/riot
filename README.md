# riot

<center>![](doc/riot.jpg)</center>

Riot is a distributed key/value system basing raft algorithm, leveldb and bitCask backend storage!!!


## Install Doc

- install grpc

```
go get -u google.golang.org/grpc
go get -a github.com/golang/protobuf/protoc-gen-go
```

- build riot

```
go build riot
```

- start a cluster

```
cd tool
bash cluster1.sh
```
