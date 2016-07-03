# Riot API

> 1. bucket just only use in boldtdb store backend
> 2. if you want to query a value with consistent, must make sure qs = 1; if you want hight performance, the qs arg equals 0 

## Bucket

### create a bucket

- `curl -v -XPOST "http://localhost:8080/riot/bucket" -d 'bucketName'`

``` http
curl -v -XPOST "http://localhost:8080/riot/bucket" -d 'student'
*   Trying 127.0.0.1...
* Connected to localhost (127.0.0.1) port 8080 (#0)
> POST /riot/bucket HTTP/1.1
> Host: localhost:8080
> User-Agent: curl/7.43.0
> Accept: */*
> Content-Length: 7
> Content-Type: application/x-www-form-urlencoded
>
* upload completely sent off: 7 out of 7 bytes
< HTTP/1.1 200 OK
< Content-Type: application/json
< Date: Sun, 03 Jul 2016 15:32:19 GMT
< Content-Length: 0
<
* Connection #0 to host localhost left intact
```

### delete a bucket

`curl -v -DELETE "http://localhost:8080/riot/bucket" -d 'bucketName'`

``` http
curl -v -XDELETE "http://localhost:8080/riot/bucket/student"
*   Trying 127.0.0.1...
* Connected to localhost (127.0.0.1) port 8080 (#0)
> DELETE /riot/bucket/student HTTP/1.1
> Host: localhost:8080
> User-Agent: curl/7.43.0
> Accept: */*
>
< HTTP/1.1 200 OK
< Content-Type: application/json
< Date: Sun, 03 Jul 2016 15:32:43 GMT
< Content-Length: 0
<
* Connection #0 to host localhost left intact
```

### get bucket info

- `curl -v "http://localhost:8080/riot/bucekt" -d 'bucketName'`

``` http
curl -v "http://localhost:8080/riot/bucket/student"
* Connected to localhost (127.0.0.1) port 8080 (#0)
> GET /riot/bucket/student HTTP/1.1
> Host: localhost:8080
> User-Agent: curl/7.43.0
> Accept: */*
>
< HTTP/1.1 200 OK
< Content-Type: application/json
< Date: Sun, 03 Jul 2016 15:33:30 GMT
< Content-Length: 201
<
* Connection #0 to host localhost left intact
{
  "BranchPageN": 0,
  "BranchOverflowN": 0,
  "LeafPageN": 0,
  "LeafOverflowN": 0,
  "KeyN": 0,
  "Depth": 1,
  "BranchAlloc": 0,
  "BranchInuse": 0,
  "LeafAlloc": 0,
  "LeafInuse": 0,
  "BucketN": 1,
  "InlineBucketN": 1,
  "InlineBucketInuse": 16
}
```

## key/value

### get value by key

- `curl -v "http://localhost:8080/riot/bucket/{key}?qs={0, 1}"`

``` http
curl -v  "http://localhost:8081/riot/key/lusi?qs=0"
*   Trying 127.0.0.1...
* Connected to localhost (127.0.0.1) port 8081 (#0)
> GET /riot/key/lusi?qs=0 HTTP/1.1
> Host: localhost:8081
> User-Agent: curl/7.43.0
> Accept: */*
> 
 
< HTTP/1.1 200 OK
< Content-Type: application/json
< Date: Sun, 03 Jul 2016 15:56:56 GMT
< Content-Length: 10
< 
* Connection #0 to host localhost left intact
{"Age":18}%                                  
```
### get value by bucket and key

be sure to exist the bucket.

- `curl -v "http://localhost:8080/riot/bucekt/{bucketName}/key/{key}?qs={0,  1}"`

``` http
curl -v  "http://localhost:8080/riot/bucket/student/key/lusi?qs=0" 
*   Trying 127.0.0.1...
* Connected to localhost (127.0.0.1) port 8080 (#0)
> GET /riot/bucket/student/key/lusi?qs=0 HTTP/1.1
> Host: localhost:8080
> User-Agent: curl/7.43.0
> Accept: */*
> 
< HTTP/1.1 200 OK
< Content-Type: application/json
< Date: Sun, 03 Jul 2016 15:47:38 GMT
< Content-Length: 11
< 
* Connection #0 to host localhost left intact
{"Age":100}%                                  
```

### set a pair of key/value

- `curl -v -XPOST "http://localhost:8080/riot/key/{key}" -d 'value'`
```http
curl -v -XPOST "http://localhost:8080/riot/key/lusi" -d '{"Age":18}' 
*   Trying 127.0.0.1...
* Connected to localhost (127.0.0.1) port 8080 (#0)
> POST /riot/key/lusi HTTP/1.1
> Host: localhost:8080
> User-Agent: curl/7.43.0
> Accept: */*
> Content-Length: 10
> Content-Type: application/x-www-form-urlencoded
> 
* upload completely sent off: 10 out of 10 bytes
< HTTP/1.1 200 OK
< Content-Type: application/json
< Date: Sun, 03 Jul 2016 15:54:27 GMT
< Content-Length: 0
< 
* Connection #0 to host localhost left intact
```

### set a pair of bucket,key/value

be sure to exist the bucket.

- `curl -v -XPOST "http://localhost:8080/riot/bucekt/{bucketName}/key/{key}" -d 'value'`

``` http
curl -v -XPOST "http://localhost:8080/riot/bucket/student/key/lusi" -d '{"Age":100}'

*   Trying 127.0.0.1...
* Connected to localhost (127.0.0.1) port 8080 (#0)
> POST /riot/bucket/student/key/lusi HTTP/1.1
> Host: localhost:8080
> User-Agent: curl/7.43.0
> Accept: */*
> Content-Length: 11
> Content-Type: application/x-www-form-urlencoded
> 
* upload completely sent off: 11 out of 11 bytes
< HTTP/1.1 200 OK
< Content-Type: application/json
< Date: Sun, 03 Jul 2016 15:45:54 GMT
< Content-Length: 0
< 
* Connection #0 to host localhost left intact
```

### delete value by key

- `curl -v -XDELETE "http://localhost:8080/riot/key/{key}"`

``` http
curl -v -XDELETE "http://localhost:8081/riot/key/lusi" 
*   Trying 127.0.0.1...
* Connected to localhost (127.0.0.1) port 8081 (#0)
> DELETE /riot/key/lusi HTTP/1.1
> Host: localhost:8081
> User-Agent: curl/7.43.0
> Accept: */*
> 
< HTTP/1.1 200 OK
< Content-Type: application/json
< Date: Sun, 03 Jul 2016 15:55:47 GMT
< Content-Length: 0
< 
* Connection #0 to host localhost left intact
```

### delete value by bucket and key 

- `curl -v -XDELETE "http://localhost:8080/riot/bucekt/{bucektName}/key/{key}"`

``` http
curl -v -XDELETE  "http://localhost:8080/riot/bucket/student/key/lusi" 
*   Trying 127.0.0.1...
* Connected to localhost (127.0.0.1) port 8080 (#0)
> DELETE /riot/bucket/student/key/lusi HTTP/1.1
> Host: localhost:8080
> User-Agent: curl/7.43.0
> Accept: */*
> 
< HTTP/1.1 200 OK
< Content-Type: application/json
< Date: Sun, 03 Jul 2016 15:49:22 GMT
< Content-Length: 0
< 
* Connection #0 to host localhost left intact
```