package client

import (
	"testing"
)

// TestClient tests client code
func TestClient(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
}
