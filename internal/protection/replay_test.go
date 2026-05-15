package protection

import (
	"testing"
	"time"
)

func TestReplayDetector_FirstSeenIsNotReplay(t *testing.T) {
	r := NewReplayDetector()
	id := int64(123)

	if replay := r.MarkAndCheck(id); replay {
		t.Fatalf("first time msgID=%d should NOT be replay", id)
	}
}

func TestReplayDetector_SecondSeenIsReplay(t *testing.T) {
	r := NewReplayDetector()
	id := int64(123)

	if replay := r.MarkAndCheck(id); replay {
		t.Fatalf("first time msgID=%d should NOT be replay", id)
	}
	if replay := r.MarkAndCheck(id); !replay {
		t.Fatalf("second time msgID=%d SHOULD be replay", id)
	}
}

func TestReplayDetector_TTLExpires(t *testing.T) {
	r := NewReplayDetector()
	id := int64(123)

	// pierwszy raz
	if replay := r.MarkAndCheck(id); replay {
		t.Fatalf("first time should NOT be replay")
	}

	// drugi raz (od razu) → replay
	if replay := r.MarkAndCheck(id); !replay {
		t.Fatalf("second time SHOULD be replay")
	}

	// symulujemy upływ czasu
	time.Sleep(replayTTL + 10*time.Millisecond)

	// po TTL → NIE replay
	if replay := r.MarkAndCheck(id); replay {
		t.Fatalf("after TTL msgID should NOT be replay")
	}
}
