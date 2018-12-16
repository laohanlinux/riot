# riot

<center>![](https://github.com/laohanlinux/riot/blob/master/doc/riot.jpg)</center>

Riot is a distributed key/value system basing raft algorithm, leveldb、 boltdb and bitCask(in fucture) backend store!!!

## Install Doc

- build riot

```shell
go build riot
```

- build riot-proxy

```shell 
go build -o riot-proxy proxy/http/bin/riot-proxy.go 
```

- start a cluster

```shell
cd tool
bash cluster.sh
```
- start proxy 

```shell 
./riot-proxy --c proxy/http/bin/cfg.toml  
```

## Api 

in doc directory

## about design detail 

[old-link](https://laohanlinux.github.io/2016/04/25/%E4%BD%BF%E7%94%A8raft%E7%AE%97%E6%B3%95%E5%BF%AB%E7%86%9F%E6%9E%84%E5%BB%BA%E4%B8%80%E4%B8%AA%E5%88%86%E5%B8%83%E5%BC%8F%E7%9A%84key-value%E7%B3%BB%E7%BB%9F/)
