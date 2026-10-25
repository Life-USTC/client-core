package pagination

import (
	"encoding/json"
	"os"
	"testing"
)

func TestProgressAgainstSharedFixtures(t *testing.T) {
	data, err := os.ReadFile("testdata/progress_cases.json")
	if err != nil {
		t.Fatal(err)
	}
	var fixtures []struct {
		Name      string
		FirstPage int
		Pages     []struct {
			ItemCount, NextPage, ReceivedCount int
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
			progress := NewProgress(tc.FirstPage)
			if progress.NextPage() != tc.FirstPage || progress.ReceivedCount() != 0 {
				t.Fatalf("initial progress = %+v", progress)
			}
			for _, page := range tc.Pages {
				progress = progress.AfterPage(page.ItemCount, page.HasMore)
				if progress.NextPage() != page.NextPage || progress.ReceivedCount() != page.ReceivedCount {
					t.Fatalf("progress = %+v, want next=%d received=%d", progress, page.NextPage, page.ReceivedCount)
				}
			}
		})
	}
}
