package playstore

import (
	"errors"
	"net/http"
	"net/url"
	"regexp"
)

var (
	appIdRegExp = regexp.MustCompile("[a-zA_Z_][\\.\\w]*")
)

var (
	ErrInvalidPackageId = errors.New("invalid package id.")
	ErrAppDoesNotExists = errors.New("the requested app wasn't found.")
)

// LookUp method always looks for the english version so it can match
// some attribute names.
func LookUp(httpGet httpGetFunc, appId string) (*App, error) {
	if !appIdRegExp.MatchString(appId) {
		return nil, ErrInvalidPackageId
	}
	res, err := doRequest(httpGet, appId, "en")
	if err != nil {
		return nil, err
	}
	document, err := NewPlayStoreDocument(res)
	if err != nil {
		return nil, err
	}
	if !isValidApp(document) {
		return nil, ErrAppDoesNotExists
	}
	return parseApp(document, "en")
}

// MultiLookUp lets you retrieve the app information translated to other
// languages. English content is always fetched.
// When the requested language is not available a new entry is added to the
// map with the key equals to the lang code and an empty value.
func MultiLookUp(httpGet httpGetFunc, appId string, languages []string) (*App, error) {
	app, err := LookUp(httpGet, appId)
	if err != nil {
		return nil, err
	}
	for _, lang := range languages {
		if lang == "en" {
			continue
		}
		res, err := doRequest(httpGet, appId, lang)
		if err != nil {
			continue
		}
		document, err := NewPlayStoreDocument(res)
		if err != nil {
			continue
		}
		app.parseDescription(document, lang)
	}
	return app, nil
}

func doRequest(httpGet httpGetFunc, appId string, lang string) (*http.Response, error) {
	url := getLookUpUrl(appId, lang)
	return httpGet(url)
}

func getLookUpUrl(appId string, lang string) *url.URL {
	query := url.Values{}
	query.Set("id", appId)
	query.Set("hl", lang)
	return &url.URL{
		Scheme:   "https",
		Host:     ENDPOINT,
		Path:     "/apps/details",
		RawQuery: query.Encode(),
	}
}
