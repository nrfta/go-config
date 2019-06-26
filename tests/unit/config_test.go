package config_test

import (
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/gobuffalo/packr"
	"github.com/neighborly/go-config"
)

type MyAppConfig struct {
	Meta config.MetaConfig

	Port            int    `mapstructure:"port"`
	BugsnagAPIKey   string `mapstructure:"bugsnag_api_key"`
	LogLevel        string `mapstructure:"log_level"`
	PrometheusPort  int    `mapstructure:"prometheus_port"`
	SegmentWriteKey string `mapstructure:"segment_write_key"`
}

var (
	testBox    = packr.NewBox(".")
	testConfig MyAppConfig
)

var _ = Describe("Test Load ", func() {
	It("It should load config file", func() {
		Expect(config.Load(testBox, &testConfig)).To(Succeed())

		Expect(testConfig.Meta.Environment).To(Equal("test"))
		Expect(testConfig.Meta.ServiceName).To(Equal("my-app"))
		Expect(testConfig.BugsnagAPIKey).To(Equal("test-key"))
		Expect(testConfig.SegmentWriteKey).To(Equal("test-write-key"))
		Expect(testConfig.LogLevel).To(Equal("error"))
		Expect(testConfig.Port).To(Equal(6002))
	})
})

var _ = Describe("Test Load With Environment Variables", func() {
	It("It should load config file with environment variables", func() {
		os.Setenv("PORT", "5001")
		os.Setenv("META_SERVICE_NAME", "service-api")

		Expect(config.Load(testBox, &testConfig)).To(Succeed())

		Expect(testConfig.Meta.Environment).To(Equal("test"))
		Expect(testConfig.Meta.ServiceName).To(Equal("service-api"))
		Expect(testConfig.BugsnagAPIKey).To(Equal("test-key"))
		Expect(testConfig.SegmentWriteKey).To(Equal("test-write-key"))
		Expect(testConfig.LogLevel).To(Equal("error"))
		Expect(testConfig.Port).To(Equal(5001))

		os.Unsetenv("PORT")
		os.Unsetenv("META_SERVICE_NAME")
	})
})
