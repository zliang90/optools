package util

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"k8s.io/klog/v2"
)

func PrettyJSON(v any) string {
	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(v); err != nil {
		return ""
	}
	return buf.String()
}

func ListContains[T comparable](items []T, v T) bool {
	for _, item := range items {
		if v == item {
			return true
		}
	}
	return false
}

type externalIP struct {
	IP string `json:"ip"`
}

func MyExternalIPAddr() string {
	c := http.Client{Timeout: 5 * time.Second}
	resp, err := c.Get("https://myexternalip.com/json")
	if err != nil {
		klog.Warning(err)
		return ""
	}
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		klog.Warning(err)
		return ""
	}
	defer resp.Body.Close()

	var v externalIP
	if err = json.Unmarshal(respBytes, &v); err != nil {
		klog.Warning(err)
		return ""
	}
	return v.IP
}
