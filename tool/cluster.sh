pkill riot 
cd ..

go build riot.go 

mv riot tool 
cd tool 

./riot -c=cfg0.toml & 
./riot -c=cfg1.toml & 
./riot -c=cfg2.toml & 
