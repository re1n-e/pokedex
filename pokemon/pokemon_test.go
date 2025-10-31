package pokemon

import (
	"testing"
)

func TestFetchPokemon(t *testing.T) {
	p := NewPokemon()

	err := p.FetchPokemon("pikachu")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if p.Name != "pikachu" {
		t.Errorf("expected name to be 'pikachu', got '%s'", p.Name)
	}

	if len(p.Abilities) == 0 {
		t.Errorf("expected at least one ability, got 0")
	}

	if p.Height == 0 {
		t.Errorf("expected non-zero height, got %d", p.Height)
	}

	if p.Weight == 0 {
		t.Errorf("expected non-zero weight, got %d", p.Weight)
	}

	if len(p.Stats) == 0 {
		t.Errorf("expected at least one stat, got 0")
	}
}
