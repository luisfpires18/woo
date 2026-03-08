package service

import "testing"

func TestIsValidKingdom(t *testing.T) {
	valid := []string{"veridor", "sylvara", "arkazia", "draxys", "nordalh"}
	for _, k := range valid {
		if !IsValidKingdom(k) {
			t.Errorf("IsValidKingdom(%q) = false, want true", k)
		}
	}

	invalid := []string{"", "moraphys", "Veridor", "zandres", "lumus", "drakanith", "foo"}
	for _, k := range invalid {
		if IsValidKingdom(k) {
			t.Errorf("IsValidKingdom(%q) = true, want false", k)
		}
	}
}
