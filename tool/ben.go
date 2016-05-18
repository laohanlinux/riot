package main

import (
	"flag"
	"sync"
	"fmt"
	"net/http"
	"github.com/laohanlinux/go-logger/logger"
	"bytes"
	"runtime"
	"io/ioutil"
)

func main() {
	var n int
	var addr string
	var qs string
	flag.IntVar(&n, "n", 0, "thread numbers")
	flag.StringVar(&addr, "addr", "localhost:8080", "")
	flag.StringVar(&qs, "qs", "1", "")
	flag.Parse()
	runtime.GOMAXPROCS(runtime.NumCPU())
	var gGroup sync.WaitGroup
	logger.Debug("Ben on Riot. n:", n)
	for i := 0; i < n; i ++ {
		gGroup.Add(1)
		logger.Info("start worker:", i)
		go func(idx int) {
			defer gGroup.Done()
			for j := 0 * idx; j < 1024; j ++ {
				httpURL := fmt.Sprintf("http://%s/riot/key/%d", addr, j)
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

	for i := 0; i < n; i ++ {
		gGroup.Add(1)
		logger.Info("start worker:", i)
		go func(idx int) {
			defer gGroup.Done()
			for j := 0 * idx; j < 1024; j ++ {
				httpURL := fmt.Sprintf("http://%s/riot/key/%d?qs=%s", addr, j, qs)
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
