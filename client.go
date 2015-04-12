package gotorrent

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"strconv"

	"github.com/optimality/gotorrent/bencoding"
)

type TrackerRequest struct {
	InfoHash   string
	PeerId     string
	Port       string
	Uploaded   int
	Downloaded int
	Left       int
	Compact    bool
	NoPeerId   bool
	Event      string
	// IP         string
	// Numwant    int
	// Key        int
	// Trackerid  int
}

var TestTrackerRequest = TrackerRequest{
	PeerId: "G001-----ABCDEFGHIJK",
	Port:   "6881",
	Event:  "started",
}

type TrackerResponse struct {
	FailureReason  string
	WarningMessage string
	Interval       int
	MinInterval    int
	TrackerId      string
	Complete       int
	Incomplete     int
	Peers          string
}

func UrlEncodeStruct(s interface{}) (m map[string][]string, err error) {
	m = map[string][]string{}
	value := reflect.ValueOf(s)
	switch value.Kind() {
	case reflect.Struct:
		for i := 0; i < value.NumField(); i++ {
			fieldName := bencoding.ToLowerCaseWithUnderscores(value.Type().Field(i).Name)
			var fieldValue []string
			field := value.Field(i)
			switch field.Kind() {
			case reflect.String:
				fieldValue = append(fieldValue, field.String())
			case reflect.Int:
				fieldValue = append(fieldValue, strconv.Itoa(int(field.Int())))
			case reflect.Bool:
				if field.Bool() {
					fieldValue = append(fieldValue, "1")
				} else {
					fieldValue = append(fieldValue, "0")
				}
			}
			if err != nil {
				return nil, err
			}
			m[fieldName] = fieldValue
		}
	default:
		return m, fmt.Errorf("Invalid argument: %v", s)
	}
	return m, nil
}

func QueryTracker(metaInfo MetaInfo) error {
	trackerRequest := TestTrackerRequest
	trackerRequest.InfoHash = metaInfo.InfoHash
	trackerRequest.Left = metaInfo.Info.Length
	urlEncodedTrackerRequest, err := UrlEncodeStruct(trackerRequest)
	if err != nil {
		return fmt.Errorf("Couldn't encode tracker request %v, err %v", trackerRequest, err)
	}

	tracker_url, err := url.Parse(metaInfo.Announce)
	if err != nil {
		return fmt.Errorf("Couldn't parse url %v", metaInfo.Announce)
	}
	values := tracker_url.Query()
	for key, value := range urlEncodedTrackerRequest {
		for _, v := range value {
			values.Add(key, v)
		}
	}
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
	var trackerResponse TrackerResponse
	bencoding.Unmarshal(string(body), &trackerResponse)
	fmt.Println(trackerResponse)
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Non-200 response: %v", resp)
	}

	return nil
}
