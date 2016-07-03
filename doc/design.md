# Riot

# http interface

```
http://localhost:8080/riot/key/xxxxxxxxxxx
```

# API

## Bucket

### create a bucket

curl -v -XPOST "http://localhost:8080/riot/bucket" -d 'bucketName'

### delete a bucket

curl -v -DELETE "http://localhost:8080/riot/bucket" -d 'bucketName'

### get bucket info

curl -v "http://localhost:8080/riot/bucekt" -d 'bucketName'

### check bucket  

curl -v "http://localhost:8080/riot/bucket" -d 'bucektName'

## key/value

### get value by key

`curl -v "http://localhost:8080/riot/bucket/{key}"`

```

```
### get value by bucket and key

curl -v "http://localhost:8080/riot/bucekt/{bucketName}/key/{key}"


### set a pair of key/value

- `curl -v -XPOST "http://localhost:8080/riot/key/{key}" -d 'value'`
```
```

### set a pair of bucket,key/value

- `curl -v -XPOST "http://localhost:8080/riot/bucekt/{bucketName}/key/{key}" -d 'value'`

```
```

### delete value by key

curl -v -XDELETE "http://localhost:8080/riot/key/{key}"

curl -v -XDELETE "http://localhost:8080/riot/bucekt/{bucektName}/key/{key}"
