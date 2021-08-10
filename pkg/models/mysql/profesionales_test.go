package mysql

import (
	"testing"
)

func TestValidname(t *testing.T) {
	t.Run("with Valid name", func(t *testing.T) {
		got := validUsername("Mathias d'Arras")
		want := true

		if got != want {
			t.Errorf("got %t want %t", got, want)
		}
	})
	t.Run("with invalid name", func(t *testing.T) {
		got := validUsername("jorge.lopez.")
		want := false

		if got != want {
			t.Errorf("got %t want %t", got, want)
		}
	})

}
