package buildinfo

import (
	"bytes"
	"log"
	"testing"
)

func TestPrintBuildInfo(t *testing.T) {
	oldFlags := log.Flags()
	defer log.SetFlags(oldFlags)
	log.SetFlags(0)

	oldOutput := log.Writer()
	defer log.SetOutput(oldOutput)

	tests := []struct {
		name     string
		version  string
		date     string
		commit   string
		expected string
	}{
		{
			name:     "all fields populated",
			version:  "1.0.0",
			date:     "2023-01-01",
			commit:   "abc123",
			expected: "Build version: 1.0.0\nBuild date: 2023-01-01\nBuild commit: abc123\n",
		},
		{
			name:     "empty version",
			version:  "",
			date:     "2023-01-01",
			commit:   "abc123",
			expected: "Build version: N/A\nBuild date: 2023-01-01\nBuild commit: abc123\n",
		},
		{
			name:     "empty date",
			version:  "1.0.0",
			date:     "",
			commit:   "abc123",
			expected: "Build version: 1.0.0\nBuild date: N/A\nBuild commit: abc123\n",
		},
		{
			name:     "empty commit",
			version:  "1.0.0",
			date:     "2023-01-01",
			commit:   "",
			expected: "Build version: 1.0.0\nBuild date: 2023-01-01\nBuild commit: N/A\n",
		},
		{
			name:     "all fields empty",
			version:  "",
			date:     "",
			commit:   "",
			expected: "Build version: N/A\nBuild date: N/A\nBuild commit: N/A\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			log.SetOutput(&buf)

			PrintBuildInfo(tt.version, tt.date, tt.commit)

			got := buf.String()
			if got != tt.expected {
				t.Errorf("PrintBuildInfo() = %q, want %q", got, tt.expected)
			}
		})
	}
}
