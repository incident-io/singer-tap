package config_test

import (
	"github.com/incident-io/singer-tap/config"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Config", func() {
	Describe("Validate", func() {
		var (
			cfg *config.Config
		)

		BeforeEach(func() {
			cfg = &config.Config{
				APIKey: "an-api-key",
			}
		})

		It("should validate", func() {
			Expect(cfg.Validate()).To(Succeed())
		})
	})
})
