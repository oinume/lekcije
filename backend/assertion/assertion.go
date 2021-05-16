package assertion

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"google.golang.org/protobuf/proto"
)

var (
	// Default options for cmp.Diff when no options specified.
	defaultCmpOptions = []cmp.Option{ //nolint
		// compare if two protocol buffer messages are equal
		cmp.Comparer(proto.Equal),
	}
)

// AssertEqual ensures that got and want are equal by cmp.Diff.
// If they are not equal, it reports failure by t.Errorf with given `message`.
func AssertEqual(t *testing.T, want, got interface{}, message string, options ...cmp.Option) {
	t.Helper()

	if len(options) == 0 {
		options = append(options, defaultCmpOptions...)
	}

	if diff := cmp.Diff(want, got, options...); diff != "" {
		if message == "" {
			t.Errorf("AssertEqual failed (-want +got):\n%s", diff)
		} else {
			t.Errorf("AssertEqual failed: %q: (-want +got):\n%s", message, diff)
		}
	}
}

// RequireEqual ensures that got and want are equal by cmp.Diff.
// If they are not equal, it reports failure by t.Fatalf with given `message`.
func RequireEqual(t *testing.T, want, got interface{}, message string, options ...cmp.Option) {
	t.Helper()

	if len(options) == 0 {
		options = append(options, defaultCmpOptions...)
	}

	if diff := cmp.Diff(want, got, options...); diff != "" {
		if message == "" {
			t.Fatalf("RequireEqual failed (-want +got):\n%s", diff)
		} else {
			t.Fatalf("RequireEqual failed: %q: (-want +got):\n%s", message, diff)
		}
	}
}
