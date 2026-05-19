package api

import (
	"testing"
	"time"
)

func TestReplayDetector_FirstSeenIsNotReplay(t *testing.T) {
	r := NewReplayDetector()
	msgID := int64(12345)

	if replay := r.MarkAndCheck(msgID); replay {
		t.Fatalf("first time msgID=%d should NOT be treated as replay", msgID)
	}
}

func TestReplayDetector_SecondSeenIsReplay(t *testing.T) {
	r := NewReplayDetector()
	msgID := int64(12345)

	if replay := r.MarkAndCheck(msgID); replay {
		t.Fatalf("first time msgID=%d should NOT be replay", msgID)
	}
	if replay := r.MarkAndCheck(msgID); !replay {
		t.Fatalf("second time msgID=%d SHOULD be replay", msgID)
	}
}

func TestReplayDetector_DifferentIDsAreIndependent(t *testing.T) {
	r := NewReplayDetector()

	id1 := int64(1)
	id2 := int64(2)

	if replay := r.MarkAndCheck(id1); replay {
		t.Fatalf("first time id1 should not be replay")
	}
	if replay := r.MarkAndCheck(id2); replay {
		t.Fatalf("first time id2 should not be replay")
	}
	if replay := r.MarkAndCheck(id1); !replay {
		t.Fatalf("second time id1 should be replay")
	}
	if replay := r.MarkAndCheck(id2); !replay {
		t.Fatalf("second time id2 should be replay")
	}
}

func TestReplayDetector_TTLExpires(t *testing.T) {
	r := NewReplayDetector()
	msgID := int64(12345)

	if replay := r.MarkAndCheck(msgID); replay {
		t.Fatalf("first time msgID should NOT be replay")
	}

	if replay := r.MarkAndCheck(msgID); !replay {
		t.Fatalf("second time msgID SHOULD be replay")
	}

	// symulacja upływu czasu
	time.Sleep(replayTTL + 10*time.Millisecond)

	if replay := r.MarkAndCheck(msgID); replay {
		t.Fatalf("after TTL msgID should NOT be replay")
	}
}
