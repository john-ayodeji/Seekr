package crawler

import (
	"reflect"
	"testing"
)

func TestParseHTML(t *testing.T) {
	tests := []struct {
		name     string
		html     string
		expected *PageData
	}{
		{
			name: "Simple page",
			html: `
				<html>
					<head>
						<title>Test Page</title>
						<meta name="description" content="Test description">
					</head>
					<body>
						<h1>Heading 1</h1>
						<h2>Heading 2</h2>
						<p>Paragraph 1</p>
						<p>Paragraph 2</p>
						<a href="/link1">Link1</a>
						<a href="https://external.com">External</a>
					</body>
				</html>
			`,
			expected: &PageData{
				Title:       "Test Page",
				Description: "Test description",
				Headings:    []string{"Heading 1", "Heading 2"},
				Paragraphs:  []string{"Paragraph 1", "Paragraph 2"},
				Links:       []string{"/link1", "https://external.com"},
			},
		},
		{
			name: "No meta description",
			html: `
				<html>
					<head><title>No Meta</title></head>
					<body>
						<h1>H1</h1>
						<p>Content</p>
						<a href="/home">Home</a>
					</body>
				</html>
			`,
			expected: &PageData{
				Title:       "No Meta",
				Description: "",
				Headings:    []string{"H1"},
				Paragraphs:  []string{"Content"},
				Links:       []string{"/home"},
			},
		},
		{
			name: "Empty page",
			html: `
				<html><head></head><body></body></html>
			`,
			expected: &PageData{
				Title:       "",
				Description: "",
				Headings:    []string{},
				Paragraphs:  []string{},
				Links:       []string{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseHTML(tt.html)
			if err != nil {
				t.Fatalf("ParseHTML() error = %v", err)
			}
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("ParseHTML() = %+v, want %+v", got, tt.expected)
			}
		})
	}
}
