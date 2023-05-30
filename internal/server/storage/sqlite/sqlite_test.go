package sqlite

import (
	"os"
	"testing"
)

func TestNewDB(t *testing.T) {
	db, err := NewDb()
	if err != nil {
		t.Errorf("error creating DB storage: %s\n", err)
		return
	}

	if db == nil {
		t.Errorf("error creating DB storage: DB storage wasn't initialized")
		return
	}

	if _, err := os.Stat("default.db"); err != nil {
		t.Errorf("error creating DB storage, %s\n", err)
		return
	}
}
