package main

import (
	"flag"
	"fmt"
)

func main() {
	var localAddr, remoteAddr string
	flag.StringVar(&localAddr, "ladrr", "localhost:8080", "local addr:port; default is localhost:8080")
	flag.StringVar(&remoteAddr, "raddrs", "", "remote addre:port ...; default is empty")
	flag.Parse()
	fmt.Println("Say Good Bye!!!")
}
