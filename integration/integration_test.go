package tap_test

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strings"

	"github.com/onsi/gomega/gexec"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Integration", Ordered, func() {
	var (
		configFile *os.File
	)

	BeforeAll(func() {
		var err error
		configFile, err = os.CreateTemp("", "config.json")
		Expect(err).ToNot(HaveOccurred())

		_, err = configFile.WriteString(fmt.Sprintf(`{"api_key": "%s"}`, os.Getenv("TEST_INCIDENT_API_KEY")))
		Expect(err).ToNot(HaveOccurred())
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

			var schema any
			data := session.Wait().Out.Contents()
			err = json.Unmarshal(data, &schema)
			Expect(err).ToNot(HaveOccurred())

			ExpectToMatchSnapshot(schema, "testdata/discover.json")
		})
	})

	Describe("Sync", func() {
		FIt("executes successfully and matches our fixtures", func() {
			cmd := exec.Command(tapPath, "--config", configFile.Name())

			session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
			Expect(err).ToNot(HaveOccurred())

			output := string(session.Wait("10s").Out.Contents())
			outputLines := strings.Split(output, "\n")

			// Collect all the records in the output stream into a list of records by stream
			// name.
			recordsByStream := map[string][]any{}
			for _, line := range outputLines {
				var object map[string]any
				err := json.Unmarshal([]byte(line), &object)
				if err != nil {
					continue
				}

				if object == nil || object["type"] != "RECORD" {
					continue
				}

				records := recordsByStream[object["stream"].(string)]
				records = append(records, object["record"])
				recordsByStream[object["stream"].(string)] = records
			}

			// Now we compare each of the record lists to those in our fixtures.
			for stream, records := range recordsByStream {
				// Sort them by ID so we don't fail to match on inconsistent ordering.
				sort.Slice(records, func(i, j int) bool {
					return records[i].(map[string]any)["id"].(string) < records[j].(map[string]any)["id"].(string)
				})

				// Now actually perform the check.
				By("stream: " + stream)
				ExpectToMatchSnapshot(records, fmt.Sprintf("testdata/sync/%s.json", stream))
			}
		})
	})
})
