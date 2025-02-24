package models

import (
	"testing"

	db "github.com/vrypan/farma/localdb"
)

func TestBasic(t *testing.T) {
	db.Open()
	key := UrlKey{
		FrameId:  1,
		UserId:   1,
		Status:   1,
		Endpoint: "endpoint1",
		Token:    "token-1",
	}
	key.Set([]byte("value1"))

	key2 := UrlKey{
		FrameId:  1,
		UserId:   1,
		Status:   1,
		Endpoint: "endpoint1",
		Token:    "token-1",
	}
	val2, err := key2.Get()
	if err != nil {
		t.Fatalf("Failed to retrieve data: %v", err)
	}
	if string(val2) != "value1" {
		t.Errorf("Expected token '%v', got '%v'", "token-1", string(val2))
	}

	key3 := UrlKey{
		FrameId:  1,
		UserId:   1,
		Status:   2,
		Endpoint: "endpoint1",
		Token:    "",
	}
	key3.Set([]byte("value2"))

	/*
		db.Set(
			[]byte("s:url:1:1:8:endpoint-2:token-2"),
			[]byte("value3"))
		db.Set(
			[]byte("s:url:1:1:8:endpoint-2:token-3"),
			[]byte("value4"))
		db.Set(
			[]byte("s:url:1:1:8:endpoint-2:"),
			[]byte(""))
	*/
}
