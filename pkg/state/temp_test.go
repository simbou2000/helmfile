package state

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestGenerateID(t *testing.T) {
	type testcase struct {
		subject string
		release ReleaseSpec
		data    any
		want    string
	}

	ids := map[string]int{}

	run := func(tc testcase) {
		t.Helper()

		t.Run(tc.subject, func(t *testing.T) {
			t.Helper()

			got, err := generateValuesID(&tc.release, tc.data)
			if err != nil {
				t.Fatalf("uenxpected error: %v", err)
			}

			if d := cmp.Diff(tc.want, got); d != "" {
				t.Fatalf("unexpected result: want (-), got (+):\n%s", d)
			}

			ids[got]++
		})
	}

	run(testcase{
		subject: "baseline",
		release: ReleaseSpec{Name: "foo", Chart: "incubator/raw"},
		want:    "foo-values-59c644b686",
	})

	run(testcase{
		subject: "different bytes content",
		release: ReleaseSpec{Name: "foo", Chart: "incubator/raw"},
		data:    []byte(`{"k":"v"}`),
		want:    "foo-values-7dcdb6d789",
	})

	run(testcase{
		subject: "different map content",
		release: ReleaseSpec{Name: "foo", Chart: "incubator/raw"},
		data:    map[string]any{"k": "v"},
		want:    "foo-values-6d6d7f4847",
	})

	run(testcase{
		subject: "different chart",
		release: ReleaseSpec{Name: "foo", Chart: "stable/envoy"},
		want:    "foo-values-7f47b755d5",
	})

	run(testcase{
		subject: "different name",
		release: ReleaseSpec{Name: "bar", Chart: "incubator/raw"},
		want:    "bar-values-7ffd78b6c7",
	})

	run(testcase{
		subject: "specific ns",
		release: ReleaseSpec{Name: "foo", Chart: "incubator/raw", Namespace: "myns"},
		want:    "myns-foo-values-76c8ff6c5",
	})

	for id, n := range ids {
		if n > 1 {
			t.Fatalf("too many occurrences of %s: %d", id, n)
		}
	}
}
