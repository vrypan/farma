package localdb

import (
	"testing"
)

func TestCreateTables(t *testing.T) {
	err := Open()
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}
	defer Close()

	err = CreateTables()
	if err != nil {
		t.Fatalf("CreateTables() failed: %v", err)
	}
}
