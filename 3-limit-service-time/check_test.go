package main

import (
	"testing"
	"time"
)

func shortTestProcess() {
	time.Sleep(3 * time.Second)
}

func longTestProcess() {
	time.Sleep(11 * time.Second)
}

func TestHandleRequest_BeginnerLevel(t *testing.T) {
	// Free user
	u1 := &User{ID: 1, IsPremium: false}

	t.Run("Short process for free user should succeed", func(t *testing.T) {
		result := HandleRequest(shortTestProcess, u1)
		if !result {
			t.Errorf("Expected short process to complete, but it was killed")
		}
	})

	t.Run("Long process for free user should be killed", func(t *testing.T) {
		result := HandleRequest(longTestProcess, u1)
		if result {
			t.Errorf("Expected long process to be killed, but it completed")
		}
	})

	// Premium user
	u2 := &User{ID: 2, IsPremium: true}

	t.Run("Long process for premium user should succeed", func(t *testing.T) {
		result := HandleRequest(longTestProcess, u2)
		if !result {
			t.Errorf("Expected premium user process to complete, but it was killed")
		}
	})
}

func TestHandleRequest_AdvancedLevel(t *testing.T) {
	u := &User{ID: 3, IsPremium: false}

	t.Run("Short process should reduce TimeUsed", func(t *testing.T) {
		result := HandleRequest(shortTestProcess, u)
		if !result {
			t.Errorf("Expected short process to succeed")
		}
		if u.TimeUsed < 3 || u.TimeUsed > 4 {
			t.Errorf("Expected TimeUsed to be around 3s, got %d", u.TimeUsed)
		}
	})

	t.Run("Another short process within quota should succeed", func(t *testing.T) {
		result := HandleRequest(shortTestProcess, u)
		if !result {
			t.Errorf("Expected second short process to succeed")
		}
	})

	t.Run("Long process beyond quota should be killed", func(t *testing.T) {
		result := HandleRequest(longTestProcess, u)
		if result {
			t.Errorf("Expected long process to be killed after quota exceeded")
		}
	})
}
