package config_test

import (
	. "chaos-galago/broker/Godeps/_workspace/src/github.com/onsi/ginkgo"
	. "chaos-galago/broker/Godeps/_workspace/src/github.com/onsi/gomega"
	. "chaos-galago/broker/config"
	"fmt"
)

var _ = Describe("#LoadConfig_#GetConfig", func() {
	It("Loads configuration from a file into a struct", func() {
		_, err := LoadConfig("fixtures/test_config.json")
		if err != nil {
			panic(fmt.Sprintf("Error loading config file [%s]...", err.Error()))
		}
		conf := GetConfig()
		Expect(conf.CatalogPath).To(Equal("test"))
		Expect(conf.DefaultProbability).To(Equal(0.4))
		Expect(conf.DefaultFrequency).To(Equal(10))
	})
})
