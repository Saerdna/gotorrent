package gotorrent

import "testing"

func TestQueryTracker(t *testing.T) {
	err := QueryTracker("http://torrent.ubuntu.com:6969/announce")
	if err != nil {
		t.Errorf("Error in query: %v", err)
	}
}
