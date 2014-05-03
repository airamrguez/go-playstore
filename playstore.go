// The playstore package lets you query the Google Play Store (Android Market).
// It basically provides two ways of querying. LookUp lets you find an app by
// id and the search method lets you find apps by search terms.
package playstore

import (
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"net/url"
)

type httpGetFunc func(url *url.URL) (*http.Response, error)

type playStoreDocument struct {
	document *goquery.Document
}

const ENDPOINT = "play.google.com/store"

// Constructs a new document wrapping goquery.document
func NewPlayStoreDocument(response *http.Response) (*playStoreDocument, error) {
	document, err := goquery.NewDocumentFromResponse(response)
	if err != nil {
		return nil, err
	}
	return &playStoreDocument{document}, nil
}

func (document *playStoreDocument) Find(selection string) (s *goquery.Selection) {
	return document.document.Find(selection)
}
