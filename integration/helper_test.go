package tap_test

import (
	"encoding/json"
	"fmt"
	"os"

	. "github.com/onsi/gomega"
)

func ExpectToMatchSnapshot(actual any, snapshotFile string) {
	// Run ginkgo with this envar to update the snapshots.
	if os.Getenv("TAP_SNAPSHOT_UPDATE") == "true" {
		fmt.Printf("Writing snapshot file %s\n", snapshotFile)

		actualJSON, err := json.MarshalIndent(actual, "", "  ")
		Expect(err).NotTo(HaveOccurred())

		err = os.WriteFile(snapshotFile, actualJSON, 0644)
		Expect(err).NotTo(HaveOccurred())
	}

	data, _ := os.ReadFile(snapshotFile)
	var expected any
	err := json.Unmarshal(data, &expected)
	Expect(err).NotTo(HaveOccurred())

	// Use google's cmp library to compare the two objects.
	Expect(actual).To(CompareEqually(expected))
}
