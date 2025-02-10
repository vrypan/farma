package fctools

import (
	"testing"
)

func Test_GetFidByUsername_vrypan(t *testing.T) {
	var username string = "vrypan"
	var expected_fid uint64 = 280

	t.Logf("Looking up fid for username=%v", username)
	hub := NewFarcasterHub()
	defer hub.Close()
	fid, err := hub.GetFidByUsername(username)
	if err != nil {
		t.Error(err)
	}
	if fid != expected_fid {
		t.Errorf("Expected fid=%v, got fid=%v", expected_fid, fid)
	}
}

func Test_likes(t *testing.T) {
	fid := uint64(280)
	hub := NewFarcasterHub()
	defer hub.Close()
	messages, err := hub.GetReactionsByFid(fid, "like", 10)
	if err != nil {
		t.Error(err)
	}
	if len(messages) != 10 {
		t.Errorf("Expected 10 likes, got %d", len(messages))
	}
	for _, m := range messages {
		t.Log(m)
	}
}
