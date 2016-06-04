package integrate_test

import (
	"io/ioutil"

	"github.com/laohanlinux/go-logger/logger"
	"github.com/laohanlinux/riot/config"

	. "github.com/laohanlinux/riot/integrate"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func reloadConfig(cfgPath string) *config.Configure {
	data, err := ioutil.ReadFile(cfgPath)
	if err != nil {
		logger.Fatal(err)
	}
	cfg, err := config.NewConfig(string(data))
	if err != nil {
		logger.Fatal(err)
	}
	return cfg
}

var _ = Describe("Riot", func() {
	var cfg *config.Configure

	// Test cfg0.tml
	var _ = Describe("cfg0.tml", func() {
		BeforeEach(func() {
			cfg = reloadConfig("tool/cfg0.tml")
		})

		Context("", func() {

		})
	})

})
