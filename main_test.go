package main

import "testing"

func TestGreet(t *testing.T) {
	if got := Greet("rollout"); got != "hello, rollout!" {
		t.Fatalf("unexpected greeting: %q", got)
	}
}

func TestGreetDefaultsToWorld(t *testing.T) {
	if got := Greet(""); got != "hello, world!" {
		t.Fatalf("unexpected default greeting: %q", got)
	}
}
