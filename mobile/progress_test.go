package mobile

import (
	"encoding/json"
	"os"
	"testing"
)

func TestPageProgressBindingAgainstSharedFixtures(t *testing.T) {
	data, err := os.ReadFile("../pagination/testdata/progress_cases.json")
	if err != nil {
		t.Fatal(err)
	}
	var fixtures []struct {
		Name      string
		FirstPage int64
		Pages     []struct {
			ItemCount, NextPage, ReceivedCount int64
			HasMore                            bool
		}
	}
	if err := json.Unmarshal(data, &fixtures); err != nil {
		t.Fatal(err)
	}
	if len(fixtures) == 0 {
		t.Fatal("missing shared progress fixtures")
	}
	for _, tc := range fixtures {
		t.Run(tc.Name, func(t *testing.T) {
			progress := NewPageProgress(tc.FirstPage)
			if progress.NextPage() != tc.FirstPage || progress.ReceivedCount() != 0 {
				t.Fatalf("initial progress = %+v", progress)
			}
			for _, page := range tc.Pages {
				previousPage, previousCount := progress.NextPage(), progress.ReceivedCount()
				next := progress.AfterPage(page.ItemCount, page.HasMore)
				if next.NextPage() != page.NextPage || next.ReceivedCount() != page.ReceivedCount {
					t.Fatalf("progress = %+v, want next=%d received=%d", next, page.NextPage, page.ReceivedCount)
				}
				if progress.NextPage() != previousPage || progress.ReceivedCount() != previousCount {
					t.Fatal("AfterPage changed the previous binding value")
				}
				progress = next
			}
		})
	}
}
