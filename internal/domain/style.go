package domain

import "fmt"

// Style identifies which template tree to emit.
//
// Adding a Style requires (1) a new constant here, (2) an entry in
// AllStyles, (3) a new embedded template tree under
// internal/adapters/templates/, and (4) a golden fixture under
// tests/golden/. Changing this list is an intentional design decision;
// the enum is closed by design so the CLI's --style flag can validate
// eagerly.
type Style string

const (
	// StyleTJ mirrors tj's convention: metatable-lazy init.lua,
	// extension-protocol table shape, plenary-busted tests.
	StyleTJ Style = "tj"

	// StylePrime mirrors ThePrimeagen's convention: singleton :setup(),
	// no default keymaps, harpoon-shape CI.
	StylePrime Style = "prime"

	// StyleOmar mirrors Omar's jedi-knights convention:
	// M.setup(opts, deps) dependency injection, detector.should_load()
	// gate, snacks-first pickers.
	StyleOmar Style = "omar"
)

// AllStyles lists every valid Style in stable display order.
var AllStyles = []Style{StyleTJ, StylePrime, StyleOmar}

// ParseStyle validates s and returns the matching Style. It returns an
// error naming the accepted values so the CLI can render a useful
// diagnostic instead of a bare "invalid".
func ParseStyle(s string) (Style, error) {
	for _, style := range AllStyles {
		if string(style) == s {
			return style, nil
		}
	}
	return "", fmt.Errorf("unknown style %q: must be one of %s", s, displayStyles())
}

func displayStyles() string {
	out := ""
	for i, s := range AllStyles {
		if i > 0 {
			out += ", "
		}
		out += string(s)
	}
	return out
}
