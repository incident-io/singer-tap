package tap_test

import (
	"fmt"

	"github.com/google/go-cmp/cmp"
	"github.com/onsi/gomega/types"
)

// CompareEqually uses Google's cmp library to calculate a diff between the expected and
// received value.
//
// It will perform semantically appropriate equality, and should be used when you're
// asking "is X meaningfully different than Y?"
//
// An example of the failure output is:
//
//	Expected
//	  *domain.Action
//	to be equal to
//	  *domain.Action
//	but got a non-empty diff:
//
//	  &domain.Action{
//	-         ID:             "01G1DP1D805D958454E11C5WY8",
//	+         ID:             "123",
//	          OrganisationID: "01G1DP1D803MG6WQPNH4M4W9CM",
//	          Organisation:   nil,
//	          ... // 17 identical fields
//	  }
//
// Especially useful to compare objects that may be pointers, or may not.
func CompareEqually(expected any) types.GomegaMatcher {
	return &compareEquallyMatcher{
		expected: expected,
	}
}

type compareEquallyMatcher struct {
	expected any
	diff     string
}

func (matcher *compareEquallyMatcher) Match(actual any) (success bool, err error) {
	matcher.diff = cmp.Diff(matcher.expected, actual)

	return matcher.diff == "", nil
}

func (matcher *compareEquallyMatcher) FailureMessage(actual any) (message string) {
	return fmt.Sprintf("Expected\n\t%T\nto be equal to\n\t%T\nbut got a non-empty diff:\n\n \033[0m %s",
		actual, matcher.expected, matcher.diff)
}

func (matcher *compareEquallyMatcher) NegatedFailureMessage(actual any) (message string) {
	return fmt.Sprintf("Expected\n\t%T\nto be non-equal to\n\t%T\nbut got a clean diff",
		actual, matcher.expected)
}
