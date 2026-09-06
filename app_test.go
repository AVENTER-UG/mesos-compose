package main

import (
	"testing"
	"time"
)

func TestReconnectBackoffGrowth(t *testing.T) {
	// Deterministic exponential growth: base 2s, doubling per attempt.
	for attempt, want := range map[int]time.Duration{
		0: 2 * time.Second,
		1: 4 * time.Second,
		2: 8 * time.Second,
		3: 16 * time.Second,
		4: 32 * time.Second,
	} {
		if got := reconnectBackoff(attempt); got != want {
			t.Errorf("reconnectBackoff(%d) = %v, want %v", attempt, got, want)
		}
	}
}

func TestReconnectBackoffCap(t *testing.T) {
	for _, attempt := range []int{5, 6, 10, 100, 1000000} {
		if got := reconnectBackoff(attempt); got != ReconnectBackoffMax {
			t.Errorf("reconnectBackoff(%d) = %v, want cap %v", attempt, got, ReconnectBackoffMax)
		}
	}
	if got := reconnectBackoff(5); got > ReconnectBackoffMax {
		t.Errorf("reconnectBackoff(5) = %v, must not exceed cap %v", got, ReconnectBackoffMax)
	}
}

func TestReconnectBackoffNegative(t *testing.T) {
	for _, attempt := range []int{-1, -42} {
		if got := reconnectBackoff(attempt); got != ReconnectBackoffBase {
			t.Errorf("reconnectBackoff(%d) = %v, want base %v", attempt, got, ReconnectBackoffBase)
		}
	}
}

func TestReconnectBackoffNonDecreasing(t *testing.T) {
	prev := reconnectBackoff(0)
	for attempt := 1; attempt <= 20; attempt++ {
		got := reconnectBackoff(attempt)
		if got < prev {
			t.Fatalf("reconnectBackoff(%d) = %v decreased below previous %v", attempt, got, prev)
		}
		prev = got
	}
}
