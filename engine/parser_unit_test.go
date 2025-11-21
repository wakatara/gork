package engine

import "testing"

func TestResolveObject(t *testing.T) {
	p := NewParser()

	tests := []struct {
		tokens []string
		want   string
	}{
		{[]string{"lamp"}, "lamp"},
		{[]string{"lantern"}, "lamp"},
		{[]string{"white", "house"}, "white-house"},
		{[]string{"kitchen", "window"}, "kitchen-window"},
		{[]string{"trophy", "case"}, "trophy-case"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			got := p.resolveObject(tt.tokens)
			if got != tt.want {
				t.Errorf("resolveObject(%v) = %q, want %q", tt.tokens, got, tt.want)
			}
		})
	}
}
