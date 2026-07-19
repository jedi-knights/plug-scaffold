package domain

import (
	"strings"
	"testing"
)

func TestParseStyle_Accepts(t *testing.T) {
	for _, want := range AllStyles {
		t.Run(string(want), func(t *testing.T) {
			got, err := ParseStyle(string(want))
			if err != nil {
				t.Fatalf("ParseStyle(%q): unexpected error %v", want, err)
			}
			if got != want {
				t.Fatalf("ParseStyle(%q) = %q, want %q", want, got, want)
			}
		})
	}
}

func TestParseStyle_RejectsUnknown(t *testing.T) {
	_, err := ParseStyle("neovide")
	if err == nil {
		t.Fatal("ParseStyle(\"neovide\"): expected error, got nil")
	}
	msg := err.Error()
	for _, s := range AllStyles {
		if !strings.Contains(msg, string(s)) {
			t.Errorf("error message %q missing accepted style %q", msg, s)
		}
	}
}

func TestParseStyle_RejectsEmpty(t *testing.T) {
	if _, err := ParseStyle(""); err == nil {
		t.Fatal("ParseStyle(\"\"): expected error, got nil")
	}
}

func TestAllStyles_StableOrder(t *testing.T) {
	// Golden-file tests and CLI help text both depend on this ordering.
	// Changing it is a breaking change to user-visible output.
	want := []Style{StyleTJ, StylePrime, StyleOmar}
	if len(AllStyles) != len(want) {
		t.Fatalf("AllStyles length = %d, want %d", len(AllStyles), len(want))
	}
	for i, s := range want {
		if AllStyles[i] != s {
			t.Errorf("AllStyles[%d] = %q, want %q", i, AllStyles[i], s)
		}
	}
}
