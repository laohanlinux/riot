// +build windows

package platform

import (
	"os"

	"github.com/laohanlinux/go-logger/logger"
)

func RegistSignal(sig ...os.Signal) {
	logger.Warn("the os platform is windows, can not handler the RegistSignal function.")
}
