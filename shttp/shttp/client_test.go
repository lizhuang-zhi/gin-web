package shttp

import (
	"context"
	"net/http"
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	// c := NewClient(WithTimeout(1 * time.Minute))
	c := NewClient(WithTimeout(20*time.Second), WithMaxIdleConnsPerHost(200))

	t.Log(c.Client.Timeout)
	t.Log(c.Client.Transport.(*http.Transport).MaxIdleConnsPerHost)
}

func TestNewHttRequest(t *testing.T) {
	req := NewHttpRequest(context.Background(), "http://www.baidu.com", "GET", nil)
	resp, err := req.Do()
	if err != nil {
		t.Error(err)
	}

	t.Log(resp)
}
