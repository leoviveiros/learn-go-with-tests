package main

import "testing"

func TestHello(t *testing.T) {
	assertMessage := func(t testing.TB, got, want string) {
		t.Helper()

		if got != want {
			t.Errorf("got %q want %q", got, want)
		}
	}

	t.Run("saying Hello to people", func(t *testing.T) {
		got := Hello("Leonardo", "")
		want := "Hello, Leonardo"

		assertMessage(t, got, want)
	})

	t.Run("say 'Hello, wold' when name is empty", func(t *testing.T) {
		got := Hello("", "")
		want := "Hello, world"

		assertMessage(t, got, want)
	})

	t.Run("in Spanish", func(t *testing.T) {
		got := Hello("Leonardo", "Spanish")
		want := "Hola, Leonardo"

		assertMessage(t, got, want)
	})

	t.Run("in French", func(t *testing.T) {
		got := Hello("Leonardo", "French")
		want := "Bonjour, Leonardo"

		assertMessage(t, got, want)
	})
}