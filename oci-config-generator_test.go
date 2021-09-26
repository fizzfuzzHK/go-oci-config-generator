package main

import "testing"

func assertCollectMessage(t *testing.T, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("got %q want %q", got, want)
	}
}
func TestGetHomeDir(t *testing.T) {

}
