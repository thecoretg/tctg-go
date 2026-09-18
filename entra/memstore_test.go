package entra

import (
	"testing"
	"time"
)

func TestMemoryStoreState(t *testing.T) {
	ctx := t.Context()
	m := NewMemoryStore()
	future := time.Now().Add(time.Minute)

	if err := m.PutState(ctx, "k", FlowState{State: "s", ExpiresAt: future}); err != nil {
		t.Fatal(err)
	}
	s, ok, err := m.TakeState(ctx, "k")
	if err != nil || !ok || s.State != "s" {
		t.Fatalf("take: %+v %v %v", s, ok, err)
	}
	if _, ok, _ := m.TakeState(ctx, "k"); ok {
		t.Fatal("state survived a take")
	}

	t.Run("cap evicts soonest to expire", func(t *testing.T) {
		m := NewMemoryStore()
		m.maxStates = 2
		now := time.Now()
		_ = m.PutState(ctx, "a", FlowState{ExpiresAt: now.Add(3 * time.Minute)})
		_ = m.PutState(ctx, "b", FlowState{ExpiresAt: now.Add(1 * time.Minute)})
		_ = m.PutState(ctx, "c", FlowState{ExpiresAt: now.Add(2 * time.Minute)})
		if _, ok, _ := m.TakeState(ctx, "b"); ok {
			t.Error("b should have been evicted")
		}
		for _, k := range []string{"a", "c"} {
			if _, ok, _ := m.TakeState(ctx, k); !ok {
				t.Errorf("%s should remain", k)
			}
		}
	})

	t.Run("expired rows are swept on put", func(t *testing.T) {
		m := NewMemoryStore()
		_ = m.PutState(ctx, "old", FlowState{ExpiresAt: time.Now().Add(-time.Second)})
		_ = m.PutState(ctx, "new", FlowState{ExpiresAt: future})
		if _, ok, _ := m.TakeState(ctx, "old"); ok {
			t.Error("expired state not swept")
		}
	})
}

func TestMemoryStoreSessions(t *testing.T) {
	ctx := t.Context()
	m := NewMemoryStore()
	now := time.Now()
	later := now.Add(time.Hour)

	_ = m.PutSession(ctx, "s1", Session{UserID: "u1", CreatedAt: now, LastSeenAt: now, ExpiresAt: later})
	_ = m.PutSession(ctx, "s2", Session{UserID: "u1", CreatedAt: now, LastSeenAt: now, ExpiresAt: later})
	_ = m.PutSession(ctx, "s3", Session{UserID: "u2", CreatedAt: now, LastSeenAt: now, ExpiresAt: later})

	if err := m.TouchSession(ctx, "s1", later, later.Add(time.Hour)); err != nil {
		t.Fatal(err)
	}
	s, ok, _ := m.GetSession(ctx, "s1")
	if !ok || !s.LastSeenAt.Equal(later) || !s.ExpiresAt.Equal(later.Add(time.Hour)) {
		t.Errorf("touch: %+v", s)
	}
	if err := m.TouchSession(ctx, "missing", later, later); err != nil {
		t.Errorf("touch missing: %v", err)
	}

	_ = m.DeleteUserSessions(ctx, "u1")
	for _, k := range []string{"s1", "s2"} {
		if _, ok, _ := m.GetSession(ctx, k); ok {
			t.Errorf("%s survived DeleteUserSessions", k)
		}
	}
	if _, ok, _ := m.GetSession(ctx, "s3"); !ok {
		t.Error("other user's session deleted")
	}

	_ = m.DeleteSession(ctx, "s3")
	if _, ok, _ := m.GetSession(ctx, "s3"); ok {
		t.Error("s3 survived DeleteSession")
	}

	_ = m.PutSession(ctx, "stale", Session{UserID: "u3", ExpiresAt: now.Add(-time.Second)})
	_ = m.PutSession(ctx, "fresh", Session{UserID: "u3", ExpiresAt: later})
	if _, ok, _ := m.GetSession(ctx, "stale"); ok {
		t.Error("expired session not swept on put")
	}
}
