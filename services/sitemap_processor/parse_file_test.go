package sitemap_processor

import (
	"fmt"
	"reflect"
	"testing"
)

func TestParseSitemap(t *testing.T) {
	passCount, failCount := 0, 0

	type testCases struct {
		Name     string
		Url      string
		Expected Sitemap
	}

	tests := []testCases{
		{
			Name: "John Ayodeji",
			Url:  "https://johnayodeji.dev/sitemap.xml",
			Expected: Sitemap{
				UrlSet: []URL{
					{
						Loc:      "https://johnayodeji.dev/",
						Priority: 1.0,
					},
					{
						Loc:      "https://johnayodeji.dev/github-stats",
						Priority: 0.6,
					},
					{
						Loc:      "https://johnayodeji.dev/projects",
						Priority: 0.7,
					},
					{
						Loc:      "https://blog.johnayodeji.dev/",
						Priority: 0.7,
					},
				},
			},
		},
		{
			Name: "Seekr",
			Url:  "https://seekr.tech/sitemap.xml",
			Expected: Sitemap{
				UrlSet: []URL{
					{
						Loc:      "https://seekr.tech/",
						Priority: 1.0,
					},
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			got, err := ParseSitemap(tc.Url)
			if err != nil {
				fmt.Println(err)
				return
			}

			if !reflect.DeepEqual(tc.Expected, got) {
				failCount += 1

				fmt.Println(got)
			}

			if reflect.DeepEqual(tc.Expected, got) {
				passCount += 1
			}
		})
	}

	fmt.Printf("\nPASS: %v, FAIL: %v\n", passCount, failCount)
}
