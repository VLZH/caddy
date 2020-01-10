package httpcache

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func newTestCache(headers []string) *Cache {
	return &Cache{
		Self:    "http://localhost:9000",
		Peers:   []string{},
		Headers: headers,
	}
}

func newTestRequest(method, url string, headers map[string]string) *http.Request {
	req := httptest.NewRequest(method, url, nil)
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	return req
}

func TestBuildKey(t *testing.T) {
	for i, c := range []struct {
		instance *Cache
		request  *http.Request
		expect   string
	}{
		{
			instance: newTestCache([]string{"Accept"}),
			request:  newTestRequest("GET", "http://test.com", map[string]string{"Accept": "image/webp,*/*"}),
			expect:   "http://test.comimage/webp,*/*",
		},
		{
			instance: newTestCache([]string{"Accept"}),
			request:  newTestRequest("GET", "http://test.com", map[string]string{}),
			expect:   "http://test.com",
		},
		{
			instance: newTestCache([]string{}),
			request:  newTestRequest("GET", "http://test.com", map[string]string{"Accept": "image/webp,*/*"}),
			expect:   "http://test.com",
		},
	} {
		key := c.instance.buildKey(c.request)
		if key != c.expect {
			t.Fatalf("Test %d: Expected key: %s, but got: %s", i, c.expect, key)
		}
	}
}
