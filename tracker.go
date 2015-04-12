package gotorrent

import (
	"fmt"
	"net/http"
	"net/url"
)

const peerID = "G001-----ABCDEFGHIJK"

func QueryTracker(metaInfo MetaInfo) error {
	tracker_url, err := url.Parse(metaInfo.Announce)
	if err != nil {
		return fmt.Errorf("Couldn't parse url %v", metaInfo.Announce)
	}
	values := tracker_url.Query()
	values.Set("peer_id", peerID)
	values.Set("port", "6881")
	values.Set("info_hash", metaInfo.InfoHash)
	tracker_url.RawQuery = values.Encode()
	fmt.Printf("querying %v\n", tracker_url)
	resp, err := http.Get(tracker_url.String())
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Non-200 response: %v", resp)
	}
	fmt.Println(resp)

	return nil
}
