package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"runtime"
	"sync"

	"github.com/laohanlinux/go-logger/logger"
)

func main() {
	var n int
	var addr string
	var qs string
	var storedb string
	flag.IntVar(&n, "n", 0, "thread numbers")
	flag.StringVar(&addr, "addr", "localhost:8080", "")
	flag.StringVar(&qs, "qs", "1", "")
	flag.StringVar(&storedb, "db", "boltdb", "")
	flag.Parse()
	runtime.GOMAXPROCS(runtime.NumCPU())
	var gGroup sync.WaitGroup
	logger.Debug("Ben on Riot. n:", n)

	if storedb == "boltdb" {
		_, err := http.Post(fmt.Sprintf("http://%s/riot/bucket", addr), "application/x-www-form-urlencoded", bytes.NewReader([]byte(`a`)))
		if err != nil {
			logger.Fatal(err)
		}
	}

	for i := 0; i < n; i++ {
		gGroup.Add(1)
		logger.Info("start worker:", i)
		go func(idx int) {
			defer gGroup.Done()
			for j := 0 * idx; j < 10; j++ {
				httpURL := fmt.Sprintf("http://%s/riot/key/%d", addr, j)
				if storedb == "boltdb" {
					httpURL = fmt.Sprintf("http://%s/riot/bucket/a/key/%d", addr, j)
				}
				logger.Debug("request in:", idx, " ", httpURL)
				ioReader := bytes.NewReader([]byte(fmt.Sprintf("%d", j)))
				resp, err := http.Post(httpURL, "application/x-www-form-urlencoded", ioReader)
				if err != nil {
					logger.Error(err)
					continue
				}
				if resp.StatusCode != 200 {
					logger.Error(resp.Status, httpURL)
					continue
				}
			}
		}(i)
	}

	gGroup.Wait()

	for i := 0; i < n; i++ {
		gGroup.Add(1)
		logger.Info("start worker:", i)
		go func(idx int) {
			defer gGroup.Done()
			for j := 0 * idx; j < 1024; j++ {
				httpURL := fmt.Sprintf("http://%s/riot/key/%d?qs=%s", addr, j, qs)
				if storedb == "boltdb" {
					httpURL = fmt.Sprintf("http://%s/riot/bucket/a/key/%d?qs=%s", addr, j, qs)
				}
				logger.Debug("request in:", idx, " ", httpURL)
				resp, err := http.Get(httpURL)
				if err != nil {
					logger.Error(err)
					continue
				}
				if resp.StatusCode != 200 {
					logger.Error(resp.Status, httpURL)
					continue
				}
				b, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					logger.Error(err)
					continue
				}
				logger.Info(j, ":", string(b))
			}
		}(i)
	}

	gGroup.Wait()
}
