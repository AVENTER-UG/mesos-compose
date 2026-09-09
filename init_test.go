package main

import "testing"

func TestParseCORSAllowedOrigins(t *testing.T) {
	got := parseCORSAllowedOrigins(` "http://localhost:5173" , https://frontend.example.invalid `)
	want := []string{"http://localhost:5173", "https://frontend.example.invalid"}
	if len(got) != len(want) {
		t.Fatalf("parseCORSAllowedOrigins() = %#v, want %#v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("parseCORSAllowedOrigins()[%d] = %q, want %q", i, got[i], want[i])
		}
	}
}
