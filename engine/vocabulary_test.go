package engine

import "testing"

func TestVocabularyLookup(t *testing.T) {
	v := NewVocabulary()

	tests := []struct {
		lookup string
		want   string
	}{
		{"lamp", "lamp"},
		{"lantern", "lamp"},
		{"white house", "white-house"},
		{"white-house", "white-house"},
		{"kitchen window", "kitchen-window"},
		{"look at", "examine"},
		{"put down", "drop"},
	}

	for _, tt := range tests {
		t.Run(tt.lookup, func(t *testing.T) {
			// Try as object
			if got := v.GetObject(tt.lookup); got != "" {
				if got != tt.want {
					t.Errorf("GetObject(%q) = %q, want %q", tt.lookup, got, tt.want)
				}
				return
			}
			// Try as verb
			if got := v.GetVerb(tt.lookup); got != "" {
				if got != tt.want {
					t.Errorf("GetVerb(%q) = %q, want %q", tt.lookup, got, tt.want)
				}
				return
			}
			t.Errorf("Vocabulary doesn't contain %q", tt.lookup)
		})
	}
}
