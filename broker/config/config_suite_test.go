package config_test

import (
	. "chaos-galago/broker/Godeps/_workspace/src/github.com/onsi/ginkgo"
	. "chaos-galago/broker/Godeps/_workspace/src/github.com/onsi/gomega"

	"testing"
)

func TestConfig(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Config test suite")
}
