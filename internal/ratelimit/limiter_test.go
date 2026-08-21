package ratelimit

import (
	"sync"
	"testing"
	"time"
)

func TestLimiter_AllowsUpToLimitThenBlocks(t *testing.T) {
	l := NewLimiter(3, time.Minute)

	for i := 0; i < 3; i++ {
		if !l.Allow("1.2.3.4") {
			t.Fatalf("attempt %d: expected allowed, got blocked", i+1)
		}
	}

	if l.Allow("1.2.3.4") {
		t.Fatal("4th attempt: expected blocked, got allowed")
	}
}

func TestLimiter_DifferentKeysAreIndependent(t *testing.T) {
	l := NewLimiter(1, time.Minute)

	if !l.Allow("1.2.3.4") {
		t.Fatal("first key: expected allowed")
	}
	if !l.Allow("5.6.7.8") {
		t.Fatal("second, different key: expected allowed regardless of the first key's state")
	}
}

func TestLimiter_ResetsAfterWindowExpires(t *testing.T) {
	l := NewLimiter(1, 50*time.Millisecond)

	if !l.Allow("1.2.3.4") {
		t.Fatal("first attempt: expected allowed")
	}
	if l.Allow("1.2.3.4") {
		t.Fatal("second attempt within window: expected blocked")
	}

	time.Sleep(60 * time.Millisecond)

	if !l.Allow("1.2.3.4") {
		t.Fatal("attempt after window expired: expected allowed")
	}
}

func TestLimiter_ConcurrentAccessIsSafe(t *testing.T) {
	l := NewLimiter(1000, time.Minute)

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			l.Allow("shared-key")
		}()
	}
	wg.Wait()
}
