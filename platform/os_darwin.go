// +build darwin

package platform

import (
	"os"
	"os/signal"

	"github.com/laohanlinux/go-logger/logger"
)

func RegistSignal(sig ...os.Signal) {
	signalChan := make(chan os.Signal)
	go func() {
		for {
			logger.Info("receive the signal: ", <-signalChan)
		}
	}()
	signal.Notify(signalChan, sig...)
}
