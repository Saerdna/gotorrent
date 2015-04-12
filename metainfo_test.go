package gotorrent

import (
	"io/ioutil"
	"testing"

	"github.com/optimality/gotorrent/bencoding"
)

func TestMetaInfo(t *testing.T) {
	testFiles := []string{
		"Plan_9_from_Outer_Space_1959_archive.torrent",
		"ubuntu-14.10-desktop-amd64.iso.torrent",
		"sample.torrent",
	}
	for _, testFile := range testFiles {
		t.Logf("Testing %v\n", testFile)
		b, err := ioutil.ReadFile("testData/" + testFile)
		if err != nil {
			t.Errorf("Unable to find testdata.")
		}
		var metaInfo MetaInfo
		err = bencoding.Unmarshal(string(b), &metaInfo)
		if err != nil {
			t.Errorf("Unable to unmarshal %v: %v", string(b), err)
		}
		t.Logf("Loaded from %v, name %v\n", metaInfo.Announce, metaInfo.Info.Name)
	}
}
