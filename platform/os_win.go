// +build windows

package platform

import (
	"os"

	"github.com/laohanlinux/go-logger/logger"
)

// RegistSignal for listening signals
func RegistSignal(sig ...os.Signal) {
	logger.Warn("the os platform is windows, can not handler the RegistSignal function.")
}
