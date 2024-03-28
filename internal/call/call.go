package call

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func GetChannel(host string) ([]string, error) {
	url := fmt.Sprintf("http://%s/channel/all", host)
	rsp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer rsp.Body.Close()
	var result []string
	if err := json.NewDecoder(rsp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return result, nil
}

func Report(data io.Reader, host, channel string) error {
	url := fmt.Sprintf("http://%s/report?channel=%s", host, channel)
	rsp, err := http.Post(url, "application/x-www-form-urlencoded", data)
	if err != nil {
		return err
	}
	defer rsp.Body.Close()
	return nil
}
