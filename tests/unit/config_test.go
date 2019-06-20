package config_test

import (
	"github.com/gobuffalo/packr"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/neighborly/go-config"
)

type DemandAPI struct {
	Meta config.MetaConfig

	Port            int    `mapstructure:"port"`
	BugsnagAPIKey   string `mapstructure:"bugsnag_api_key"`
	LogLevel        string `mapstructure:"log_level"`
	PrometheusPort  int    `mapstructure:"prometheus_port"`
	SegmentWriteKey string `mapstructure:"segment_write_key"`
}

var (
	testBox    = packr.NewBox(".")
	testConfig DemandAPI
)

var _ = Describe("Test Load ", func() {
	It("It should load config file", func() {

		Expect(config.Load(testBox, &testConfig)).To(Succeed())

	})

})

var _ = Describe("Test GetMetaConfig ", func() {
	It("It should get the  meta values of the  config ", func() {
		Expect(config.GetMetaConfig(&testConfig)).To(Not(BeNil()))
	})

})
