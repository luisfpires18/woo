package service

import "testing"

func TestIsValidKingdom(t *testing.T) {
	valid := []string{"veridor", "sylvara", "arkazia", "draxys", "nordalh", "zandres", "lumus"}
	for _, k := range valid {
		if !IsValidKingdom(k) {
			t.Errorf("IsValidKingdom(%q) = false, want true", k)
		}
	}

	invalid := []string{"", "moraphys", "Veridor", "drakanith", "foo"}
	for _, k := range invalid {
		if IsValidKingdom(k) {
			t.Errorf("IsValidKingdom(%q) = true, want false", k)
		}
	}
}
