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

func TestCategory(t *testing.T) {
	app, err := LookUp(httpGet, "com.google.android.youtube")
	if err != nil {
		t.Logf("Error: %s.", err.Error())
	} else {
		if app.Category != "Media & Video" {
			t.Error("Expected app category to be the same.")
		}
	}
}

func TestDeveloperName(t *testing.T) {
	app, err := LookUp(httpGet, "com.google.android.youtube")
	if err != nil {
		t.Logf("Error: %s.", err.Error())
	} else {
		if app.Developer.Name != "Google Inc." {
			t.Error("Expected app category to be the same.")
		}
	}
}

func httpGet(url *url.URL) (*http.Response, error) {
	return http.Get(url.String())
}
