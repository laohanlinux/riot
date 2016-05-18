pid=`ps -ef | grep riot | awk '{print $2}'`
kill `echo $pid`

cd ..
go build riot.go 

mv riot tool 
cd tool

rm -fr raft0 raft1 raft2 raft3 raft4

./riot -c=cfg0.toml & 

sleep 2

./riot -c=cfg1.toml -join="127.0.0.1:8080" & 
sleep 2

./riot -c=cfg2.toml -join="127.0.0.1:8080" & 
