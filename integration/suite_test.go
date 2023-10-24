package tap_test

import (
	"testing"

	"github.com/onsi/gomega/gexec"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestSuite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "integration")
}

var tapPath string

// Before we run any tests, we need to compile our tap so we can execute the binary.
var _ = BeforeSuite(func() {
	var err error
	tapPath, err = gexec.Build("github.com/incident-io/singer-tap/cmd/tap-incident")
	Expect(err).ToNot(HaveOccurred())
})

var _ = AfterSuite(func() {
	gexec.CleanupBuildArtifacts()
})
