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

curl -v "http://localhost:8080/bucket" -d 'bucektName'

## key/value
