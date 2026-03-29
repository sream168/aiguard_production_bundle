package strutil

import "testing"

func TestFirstNonEmpty(t *testing.T) {
	tests := []struct {
		name   string
		input  []string
		expect string
	}{
		{"all empty", []string{"", "  ", ""}, ""},
		{"first non-empty", []string{"", "hello", "world"}, "hello"},
		{"trims spaces", []string{"  ", " hi "}, "hi"},
		{"single value", []string{"only"}, "only"},
		{"no values", []string{}, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FirstNonEmpty(tt.input...)
			if got != tt.expect {
				t.Errorf("FirstNonEmpty(%v) = %q, want %q", tt.input, got, tt.expect)
			}
		})
	}
}
