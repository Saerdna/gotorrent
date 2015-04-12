package gotorrent

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
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
	values.Set("uploaded", "0")
	values.Set("downloaded", "0")
	values.Set("left", strconv.Itoa(metaInfo.Info.Length))
	values.Set("compact", "0")
	values.Set("no_peer_id", "0")
	// values.Set("event", "started")
	tracker_url.RawQuery = values.Encode()

	fmt.Printf("querying %v\n", tracker_url)
	resp, err := http.Get(tracker_url.String())
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	fmt.Println(resp)
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	fmt.Println(string(body))
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Non-200 response: %v", resp)
	}

	return nil
}
