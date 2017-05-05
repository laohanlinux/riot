// +build linux

package platform

import (
	"os"
	"os/signal"

	log "github.com/laohanlinux/utils/gokitlog"
)

// RegistSignal for listening signals
func RegistSignal(sig ...os.Signal) {
	signalChan := make(chan os.Signal)
	go func() {
		for {
			log.Info("receive the signal: ", <-signalChan)
		}
	}()
	signal.Notify(signalChan, sig...)
}
