// +build windows

package platform

import (
	"os"

	log "github.com/laohanlinux/utils/gokitlog"
)

// RegistSignal for listening signals
func RegistSignal(sig ...os.Signal) {
	log.Warn("the os platform is windows, can not handler the RegistSignal function.")
}
