package navigation

import (
	"encoding/json"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHeadersMarshalJSON(t *testing.T) {
	headers := Headers{
		"Content-Type": "application/json",
		"X-CUSTOM":     "value",
	}

	serialized, err := headers.MarshalJSON()
	require.NoError(t, err)

	var decoded map[string]string
	require.NoError(t, json.Unmarshal(serialized, &decoded))

	expected := map[string]string{
		"content-type": "application/json",
		"x-custom":     "value",
	}
	assert.Equal(t, expected, decoded)

	_, stillExists := headers["Content-Type"]
	assert.True(t, stillExists, "original map should retain mixed-case keys")
}

func TestResponseAbsoluteURL(t *testing.T) {
	baseURL, err := url.Parse("https://example.com/path/index.html?x=1#anchor")
	require.NoError(t, err)

	response := Response{
		Resp: &http.Response{Request: &http.Request{URL: baseURL}},
	}

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "fragment only",
			input:    "#section",
			expected: "",
		},
		{
			name:     "relative path removes fragment",
			input:    "../other/page.html#frag",
			expected: "https://example.com/other/page.html",
		},
		{
			name:     "protocol relative path",
			input:    "//cdn.example.com/script.js",
			expected: "https://cdn.example.com/script.js",
		},
		{
			name:     "invalid escape returns empty",
			input:    "%zz",
			expected: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := response.AbsoluteURL(tc.input)
			assert.Equal(t, tc.expected, got)
		})
	}
}

func TestResponseIsRedirect(t *testing.T) {
	assert.True(t, Response{StatusCode: 302}.IsRedirect())
	assert.False(t, Response{StatusCode: 404}.IsRedirect())
}

func TestRequestURL(t *testing.T) {
	tests := []struct {
		name     string
		request  Request
		expected string
	}{
		{
			name: "GET returns URL",
			request: Request{
				Method: http.MethodGet,
				URL:    "https://example.com/get",
			},
			expected: "https://example.com/get",
		},
		{
			name: "POST combines URL and body",
			request: Request{
				Method: http.MethodPost,
				URL:    "https://example.com/post",
				Body:   "payload=true",
			},
			expected: "https://example.com/post:payload=true",
		},
		{
			name: "unsupported method returns empty",
			request: Request{
				Method: http.MethodPut,
				URL:    "https://example.com/put",
			},
			expected: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.request.RequestURL()
			assert.Equal(t, tc.expected, got)
		})
	}
}

func TestNewNavigationRequestURLFromResponse(t *testing.T) {
	baseURL, err := url.Parse("https://example.org/base/index.html")
	require.NoError(t, err)

	resp := &Response{
		Resp:         &http.Response{Request: &http.Request{URL: baseURL}},
		Depth:        3,
		RootHostname: "example.org",
	}

	request := NewNavigationRequestURLFromResponse("/new/path", "source", "tag", "href", resp)

	assert.Equal(t, http.MethodGet, request.Method)
	assert.Equal(t, "https://example.org/new/path", request.URL)
	assert.Equal(t, 3, request.Depth)
	assert.Equal(t, "example.org", request.RootHostname)
	assert.Equal(t, "source", request.Source)
	assert.Equal(t, "tag", request.Tag)
	assert.Equal(t, "href", request.Attribute)
}
