package httpcache

import (
	"crypto/md5"
	"fmt"
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

func getMD5String(key string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(key)))
}

func TestBuildKey(t *testing.T) {
	for i, c := range []struct {
		instance *Cache
		request  *http.Request
		expect   string
	}{
		{
			instance: newTestCache([]string{"Accept", "Accept-Encoding"}),
			request: newTestRequest("GET", "http://test.com", map[string]string{
				"Accept":          "image/webp,*/*",
				"Accept-Encoding": "gzip, deflate, br",
			}),
			expect: getMD5String("http://test.comimage/webp,*/*gzip, deflate, br"),
		},
		{
			instance: newTestCache([]string{"Accept-Encoding", "Accept"}),
			request: newTestRequest("GET", "http://test.com", map[string]string{
				"Accept":          "image/webp,*/*",
				"Accept-Encoding": "gzip, deflate, br",
			}),
			expect: getMD5String("http://test.comgzip, deflate, brimage/webp,*/*"),
		},
		{
			instance: newTestCache([]string{"Accept"}),
			request:  newTestRequest("GET", "http://test.com", map[string]string{}),
			expect:   getMD5String("http://test.com"),
		},
		{
			instance: newTestCache([]string{}),
			request:  newTestRequest("GET", "http://test.com", map[string]string{"Accept": "image/webp,*/*"}),
			expect:   getMD5String("http://test.com"),
		},
	} {
		key := c.instance.buildKey(c.request)
		if key != c.expect {
			t.Fatalf("Test %d: Expected key: %s, but got: %s", i, c.expect, key)
		}
	}
}
