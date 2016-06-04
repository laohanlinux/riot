package integrate_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestIntegrate(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Integrate Suite")
}
