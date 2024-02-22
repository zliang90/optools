package util

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/pkg/errors"
)

var defaultTimeout = 15 * time.Second

// RequestJSON request json api, and unmarshal to a specific type
// eg:
//
// var users []User
// code, err := RequestJONS("POST", "http://api.sv.com/users", body, &users)
// if err != nil {
//   return err
// }
//

func RequestJSON(method, url string, headers map[string]string, body, to any) (statusCode int, err error) {
	var (
		reqBytes []byte
		br       io.Reader
	)

	if body != nil {
		reqBytes, err = json.Marshal(body)
		if err != nil {
			return 0, err
		}
		br = bytes.NewReader(reqBytes)
	}

	var req *http.Request
	req, err = http.NewRequest(method, url, br)
	if err != nil {
		return 0, err
	}

	req.Header.Set("Accept", "application/json")
	if ListContains([]string{"POST", "PUT", "PATCH"}, method) && body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if headers != nil {
		for k, v := range headers {
			req.Header.Add(k, v)
		}
	}

	var resp *http.Response
	resp, err = (&http.Client{Timeout: defaultTimeout}).Do(req)
	if err != nil {
		return 0, err
	}
	statusCode = resp.StatusCode

	var respBytes []byte
	respBytes, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	if statusCode != http.StatusOK {
		err = errors.Errorf("unexpected response, http code: %d, body: %q",
			resp.StatusCode, string(respBytes))
		return
	}

	if to != nil {
		if err = json.Unmarshal(respBytes, to); err != nil {
			return
		}
	}

	return
}
