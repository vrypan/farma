package models

import (
	"bytes"
	"crypto/ed25519"
	"crypto/rand"
	"testing"
	//db "github.com/vrypan/farma/localdb"
)

func Test_Pubkey_1(t *testing.T) {
	pubKey, privKey, err := ed25519.GenerateKey(rand.Reader)
	t.Logf("Generated public key: %x", pubKey)
	t.Logf("Generated private key: %x", privKey)
	if err != nil {
		t.Errorf("Error generating key pair: %v", err)
	}
	pk := PubKey{Key: pubKey, FrameId: "1"}
	t.Logf("DB Key is %s", pk.DbKey())
	httpHeader := pk.Encode()
	t.Logf("HTTP Header: %s", httpHeader)

	pk2 := PubKey{}
	pk2.Decode(httpHeader)

	if !bytes.Equal(pk2.Key, pk.Key) {
		t.Errorf("Expected key %x, got %x", pk.Key, pk2.Key)
	}
	if pk2.FrameId != pk.FrameId {
		t.Errorf("Expected frame id %s, got %s", pk.FrameId, pk2.FrameId)
	}

}

/*
func pubkey_cleanup(t *testing.T) int {
	keys, _, _ := db.GetKeysWithPrefix([]byte("n:id:test-"), []byte("n:id:test-"), 1000)
	for _, key := range keys {
		err := db.Delete(key)
		if err != nil {
			t.Logf("Error deleting key %s: %v\n", key, err)
		}
		//t.Logf("Deleted key %s\n", key)
	}
	return len(keys)
}

func TestNotification_Save_Load(t *testing.T) {
	db.Open()
	defer db.Close()
	cleanup(t)
	n := NewNotification(
		"test-0001",
		"Test title",
		"Test message",
		"https://link.example.com",
		"https://endpoint.example.com",
		[][]byte{},
	)
	n.Save()
	t.Logf("Saved one entry, with Id %s", n.Id)
	n2 := &Notification{}
	notificationList, err := n2.Load("test-0001")
	if err != nil {
		t.Errorf("Error loading notification: %v", err)
	}
	if len(notificationList) != 1 {
		t.Errorf("Expected 1 notification, got %d", len(notificationList))
	}
	t.Logf("Loaded one entry, with Id %s", notificationList[0].Id)
	if notificationList[0].Id != "test-0001" {
		t.Errorf("Expected notification ID 'test-0001', got %s", notificationList[0].Id)
	}
	if notificationList[0].Title != "Test title" {
		t.Errorf("Expected notification title 'Test title', got %s", notificationList[0].Title)
	}
	if notificationList[0].Message != "Test message" {
		t.Errorf("Expected notification message 'Test message', got %s", notificationList[0].Message)
	}
	if notificationList[0].Link != "https://link.example.com" {
		t.Errorf("Expected notification link 'https://link.example.com', got %s", notificationList[0].Link)
	}
	if notificationList[0].Endpoint != "https://endpoint.example.com" {
		t.Errorf("Expected notification endpoint 'https://endpoint.example.com', got %s", notificationList[0].Endpoint)
	}
	if err := db.Delete(n.PrefixBytes()); err != nil {
		t.Errorf("Error deleting notification: %v", err)
	}
	if c := cleanup(t); c != 1 {
		t.Errorf("Expected cleanup to return 1, got %d", c)
	}
}

func TestNotification_Versions(t *testing.T) {
	db.Open()
	defer db.Close()
	cleanup(t)
	n := NewNotification(
		"test-0002",
		"Test title",
		"Test message",
		"https://link.example.com",
		"https://endpoint.example.com",
		[][]byte{},
	)
	v, err := n.Save()
	if err != nil {
		t.Errorf("Error saving notification: %v", err)
	}
	if v != 0 {
		t.Errorf("Expected version 0, got %d", v)
	}
	v, err = n.Save()
	if err != nil {
		t.Errorf("Error saving notification: %v", err)
	}
	if v != 1 {
		t.Errorf("Expected version 1, got %d", v)
	}
	v, err = n.Save()
	if err != nil {
		t.Errorf("Error saving notification: %v", err)
	}
	if v != 2 {
		t.Errorf("Expected version 2, got %d", v)
	}
	t.Logf("Saved 3 entries, with Id %s", n.Id)
	n2 := &Notification{}
	notificationList, err := n2.Load("test-0002")
	if err != nil {
		t.Errorf("Error loading notification: %v", err)
	}
	if len(notificationList) != 3 {
		t.Errorf("Expected 3 notifications, got %d", len(notificationList))
	}
	t.Logf("Loaded 3 entries, with Id %s", notificationList[0].Id)
	for i := range 3 {
		if notificationList[i].Id != "test-0002" {
			t.Errorf("Expected notification ID 'test-0002', got %s", notificationList[i].Id)
		}
		if notificationList[i].Title != "Test title" {
			t.Errorf("Expected notification title 'Test title', got %s", notificationList[i].Title)
		}
		if notificationList[i].Message != "Test message" {
			t.Errorf("Expected notification message 'Test message', got %s", notificationList[i].Message)
		}
		if notificationList[i].Link != "https://link.example.com" {
			t.Errorf("Expected notification link 'https://link.example.com', got %s", notificationList[i].Link)
		}
		if notificationList[i].Endpoint != "https://endpoint.example.com" {
			t.Errorf("Expected notification endpoint 'https://endpoint.example.com', got %s", notificationList[i].Endpoint)
		}
		if notificationList[i].GetVersion() != uint64(i) {
			t.Errorf("Expected notification version to be %d, got %d", uint64(i), notificationList[i].GetVersion())
		}
	}

	if c := cleanup(t); c != 3 {
		t.Errorf("Expected cleanup to return 1, got %d", c)
	}
}

/*
func TestListAllKeys(t *testing.T) {
	db.Open()
	defer db.Close()
	keys, _, _ := db.GetKeysWithPrefix([]byte("n:test-"), []byte("n:test-"), 1000)
	for _, key := range keys {
		fmt.Printf("%s\n", key)
	}
}
*/
