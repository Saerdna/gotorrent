package gotorrent

import (
	"fmt"
	"net/http"
)

func QueryTracker(tracker_url string) error {
	resp, err := http.Get(tracker_url)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("Non-200 response: %v", resp)
	}
	fmt.Println(resp)

	return nil
}
