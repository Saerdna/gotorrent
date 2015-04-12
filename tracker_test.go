package gotorrent

import (
	"io/ioutil"
	"testing"

	"github.com/optimality/gotorrent/bencoding"
)

func TestQueryTracker(t *testing.T) {
	testFile := "ubuntu-14.10-desktop-amd64.iso.torrent"
	b, err := ioutil.ReadFile("testData/" + testFile)
	if err != nil {
		t.Errorf("Unable to find testdata.")
	}
	var metaInfo MetaInfo
	err = bencoding.Unmarshal(string(b), &metaInfo)
	if err != nil {
		t.Errorf("Unable to unmarshal %v: %v", string(b), err)
	}
	err = QueryTracker(metaInfo)
	if err != nil {
		t.Errorf("Error in query: %v", err)
	}
}
