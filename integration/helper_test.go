package tap_test

import (
	"encoding/json"
	"fmt"
	"os"

	. "github.com/onsi/gomega"
)

func ExpectJSONToMatchSnapshot(actualJSON []byte, snapshotFile string) {
	// Lets indent the json string we have so that it's easier to read.
	// We'll need to unmarshal then remarshal to do this.
	var actual interface{}
	err := json.Unmarshal(actualJSON, &actual)
	Expect(err).NotTo(HaveOccurred())

	actualJSON, err = json.MarshalIndent(actual, "", "  ")
	Expect(err).NotTo(HaveOccurred())

	// Run ginkgo with this envar to update the snapshots.
	if os.Getenv("TAP_SNAPSHOT_UPDATE") == "true" {
		fmt.Printf("Writing snapshot file %s\n", snapshotFile)

		err = os.WriteFile(snapshotFile, actualJSON, 0644)
		Expect(err).NotTo(HaveOccurred())
	}

	data, _ := os.ReadFile(snapshotFile)
	Expect(actualJSON).To(MatchJSON(data))
}
