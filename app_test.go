package playstore

import (
	"net/http"
	"net/url"
	"testing"
)

func TestTitle(t *testing.T) {
	app, err := LookUp(httpGet, "com.google.android.youtube")
	if err != nil {
		t.Logf("Error: %s.", err.Error())
	} else {
		if app.Title != "YouTube" {
			t.Error("Expected app title to be the same.")
		}
	}
}

func httpGet(url *url.URL) (*http.Response, error) {
	return http.Get(url.String())
}
