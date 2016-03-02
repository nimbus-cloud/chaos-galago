package utils_test

import (
	. "chaos-galago/processor/Godeps/_workspace/src/github.com/onsi/ginkgo"
	. "chaos-galago/processor/Godeps/_workspace/src/github.com/onsi/gomega"
	"testing"
)

func TestUtils(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Utility test suite")
}
