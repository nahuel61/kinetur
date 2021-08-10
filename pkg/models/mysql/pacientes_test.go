package mysql

import (
	"testing"
)

func TestValidEmail(t *testing.T) {
	t.Run("with Valid email", func(t *testing.T) {
		got := validEmail("pepe@gmail.com")
		want := true
		if got != want {
			t.Errorf("got %t want %t", got, want)
		}
	})

	t.Run("with invalid utf8 character", func(t *testing.T) {
		got := validEmail("ñsalazar@fie.edu.ar")
		want := false
		if got != want {
			t.Errorf("got %t want %t", got, want)
		}
	})

	t.Run("with No Domain", func(t *testing.T) {
		got := validEmail("nsalazar")
		want := false
		if got != want {
			t.Errorf("got %t want %t", got, want)
		}
	})

	t.Run("with no username", func(t *testing.T) {
		got := validEmail("@hotmail.com.ar")
		want := false
		if got != want {
			t.Errorf("got %t want %t", got, want)
		}
	})
}

func TestValidDNI(t *testing.T) {
	t.Run("with Valid DNI", func(t *testing.T) {
		got := validDNI("12345678")
		want := true
		if got != want {
			t.Errorf("got %t want %t", got, want)
		}
	})
	t.Run("more than 8 digits DNI", func(t *testing.T) {
		got := validDNI("123456789")
		want := false
		if got != want {
			t.Errorf("got %t want %t", got, want)
		}
	})
	t.Run("less than 8 digits DNI", func(t *testing.T) {
		got := validDNI("123456")
		want := false
		if got != want {
			t.Errorf("got %t want %t", got, want)
		}
	})

}
func TestValidPasswd(t *testing.T) {
	t.Run("with Valid Password", func(t *testing.T) {
		got := validPassword("qwerty")
		want := true
		if got != want {
			t.Errorf("got %t want %t", got, want)
		}
	})
	t.Run("Short password", func(t *testing.T) {
		got := validPassword("qwert")
		want := false
		if got != want {
			t.Errorf("got %t want %t", got, want)
		}
	})
}
