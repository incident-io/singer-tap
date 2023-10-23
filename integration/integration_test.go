//go:build integration
// +build integration

package tap_test

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

var _ = Describe("Integration", Ordered, func() {
	var configFile *os.File

	BeforeAll(func() {
		var err error
		configFile, err = os.CreateTemp("", "config.json")
		if err != nil {
			log.Fatal(err)
		}
		configFile.WriteString(fmt.Sprintf(`{"api_key": "%s"}`, os.Getenv("TEST_INCIDENT_API_KEY")))
	})

	AfterAll(func() {
		os.Remove(configFile.Name())
	})

	Describe("Discover", func() {
		It("runs without erroring", func() {
			cmd := exec.Command(tapPath, "--discover", "--config", configFile.Name())
			session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
			Expect(err).ToNot(HaveOccurred())

			Eventually(session).Should(gexec.Exit(0))
		})

		It("returns the full schema as expected", func() {
			cmd := exec.Command(tapPath, "--discover", "--config", configFile.Name())
			session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
			Expect(err).ToNot(HaveOccurred())

			ExpectJSONToMatchSnapshot(session.Wait().Out.Contents(), "testdata/discover.json")
		})
	})
})
