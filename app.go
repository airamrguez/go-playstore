package playstore

import (
	"github.com/PuerkitoBio/goquery"
	"regexp"
	"strconv"
	"strings"
)

// DeveloperSlug represents all the developer information contained into
// google's search results.
type DeveloperSlug struct {
	Name string `json:"name"`
}

// AppSlug contains all the information related to the app contained into
// google's search results.
type AppSlug struct {
	Title         string        `json:"title"`
	Icon          string        `json:"icon"`
	AverageRating float64       `json:"average_rating"`
	Developer     DeveloperSlug `json:"developer"`
	Price         string        `json:"price"`
}

// Creates a new AppSlug struct to store the minimum information of an app
// returned from a search method.
func NewAppSlug() *AppSlug {
	app := &AppSlug{}
	app.Developer = DeveloperSlug{}
	return app
}

// This methods collects all the fields returned from a search that are
// contained into what google call cards in the Google Play Store.
func parseAppSlug(sel *goquery.Selection) (*AppSlug, error) {
	app := NewAppSlug()
	app.parseTitle(sel)
	app.parseIcon(sel)
	app.parseAverageRating(sel)
	app.parsePrice(sel)
	return app, nil
}

// Parses the application title returned from a search.
func (app *AppSlug) parseTitle(sel *goquery.Selection) {
	app.Title = strings.TrimSpace(sel.Find(`.title`).First().Text())
}

// Parses the application icon returned from a search.
func (app *AppSlug) parseIcon(sel *goquery.Selection) {
	if icon, ok := sel.Find(`.cover-image`).First().Attr("src"); ok {
		app.Icon = icon
	}
}

// Parses the application developer name returned from a search.
func (app *AppSlug) parseDeveloper(sel *goquery.Selection) {
	app.Developer.Name = strings.TrimSpace(sel.Find(`.subtitle`).First().Text())
}

// Parses the application average rating returned from a search. In case there's
// no width attribute let's suppose that google is omitting the with attribute.
func (app *AppSlug) parseAverageRating(sel *goquery.Selection) {
	width, ok := sel.Find(`.current-rating`).First().Attr("style")
	if !ok {
		app.AverageRating = 0
		return
	}
	re := regexp.MustCompile(`width:\s*([0-9.]+)%`)
	app.AverageRating = ParseFloat(re.ReplaceAllString(width, "$1"))
}

// Parses the application price returned from a search.
func (app *AppSlug) parsePrice(sel *goquery.Selection) {
	app.Price = strings.TrimSpace(sel.Find(`.price-container button.price`).First().Text())
}

// Extends the developer slug information by adding the developer mail and website.
type Developer struct {
	DeveloperSlug
	Email   string `json:"email"`
	Website string `json:"website"`
}

// App struct represents all the information contained into an AppSlug and the
// rest of the fields returned from a lookup.
type App struct {
	AppSlug
	Category         string            `json:"category"`
	OffersInApp      bool              `json:"offers_inapp"`
	Rating           map[string]int64  `json:"rating"`
	Reviews          int64             `json:"reviews"`
	ScreenshotUrls   []string          `json:"screenshot_urls"`
	Updated          string            `json:"updated"`
	Version          string            `json:"version"`
	Size             string            `json:"size"`
	RequiresAndroid  string            `json:"requires_android"`
	ContentRating    string            `json:"content_rating"`
	Installs         string            `json:"installs"`
	PlainDescription map[string]string `json:"plain_description"`
	HtmlDescription  map[string]string `json:"html_description"`
	Developer        Developer         `json:"developer"`
}

// Creates a new App.
func NewApp() *App {
	app := &App{}
	app.Rating = make(map[string]int64, 5)
	app.PlainDescription = map[string]string{}
	app.HtmlDescription = map[string]string{}
	app.Developer = Developer{}
	return app
}

// Looks through the document extracting all the app information.
func parseApp(document *playStoreDocument, lang string) (*App, error) {
	app := NewApp()
	app.parseTitle(document)
	app.parseIcon(document)
	app.parseOffersInApp(document)
	app.parsePrice(document)
	app.parseAverageRating(document)
	app.parseReviews(document)
	app.parseMetaInfo(document)
	app.parseDescription(document, lang)
	app.parseScreenshotUrls(document)
	app.parseDeveloperName(document)
	app.parseCategory(document)
	app.parseRating(document)
	return app, nil
}

// Parses the application title returned from a lookup.
func (app *App) parseTitle(document *playStoreDocument) {
	title := document.Find(`.document-title`).Children().First().Text()
	app.Title = strings.TrimSpace(title)
}

// Parses the application icon returned from a lookup.
func (app *App) parseIcon(document *playStoreDocument) {
	if icon, ok := document.Find(`.details-info .cover-image`).First().Attr("src"); ok {
		app.Icon = icon
	}
}

// Parses the user average rating.
func (app *App) parseAverageRating(document *playStoreDocument) {
	avgRatingText := strings.TrimSpace(document.Find(`.score-container .score`).Text())
	if avgRating, err := strconv.ParseFloat(avgRatingText, 32); err == nil {
		app.AverageRating = avgRating
		return
	}
	app.AverageRating = -1.
}

// Parses the amount of reviews for the app
func (app *App) parseReviews(document *playStoreDocument) {
	reviewsText := strings.TrimSpace(document.Find(`.score-container .reviews-num`).Text())
	if reviews, err := strconv.ParseInt(reviewsText, 10, 32); err == nil {
		app.Reviews = reviews
	}
	app.Reviews = -1
}

// A more detailed view of the app is inside a container with the class
// .meta-info. Here we collect the last time the app was updated, it's current
// size, the number of installations, the minium Android required version,
// content rating and the e-mail and the website of the developer.
func (app *App) parseMetaInfo(document *playStoreDocument) {
	document.Find(`.meta-info`).Each(func(i int, sel *goquery.Selection) {
		title := strings.TrimSpace(sel.Find(`.title`).Text())
		switch title {
		case "Updated":
			app.parseUpdated(sel)
		case "Size":
			app.parseSize(sel)
		case "Installs":
			app.parseInstalls(sel)
		case "Current Version":
			app.parseVersion(sel)
		case "Requires Android":
			app.parseRequiresAndroid(sel)
		case "Content Rating":
			app.parseContentRating(sel)
		case "Contact Developer":
			app.parseContactDeveloper(sel)
		}
	})
}

// Parses the app description an stores it in plain text and html.
func (app *App) parseDescription(document *playStoreDocument, lang string) {
	htmlLang, ok := document.Find("html").Attr("lang")
	if !ok || htmlLang != lang {
		return
	}
	description := document.Find(`.details-section.description .id-app-orig-desc`)
	app.HtmlDescription[lang], _ = description.Html()
	app.PlainDescription[lang] = description.Text()
}

// Parses every src attribute from all image elements tagged with screenshot.
func (app *App) parseScreenshotUrls(document *playStoreDocument) {
	document.Find(`img.screenshot`).Each(func(i int, sel *goquery.Selection) {
		screenshotUrl, found := sel.Attr("src")
		if !found {
			return
		}
		app.ScreenshotUrls = append(app.ScreenshotUrls, screenshotUrl)
	})
}

// Parses the developer name returned from a lookup.
func (app *App) parseDeveloperName(document *playStoreDocument) {
	app.Category = strings.TrimSpace(document.Find(`.category`).Find(`[itemprop="genre"]`).First().Text())
}

// Parses the category returned from a lookup.
func (app *App) parseCategory(document *playStoreDocument) {
	app.Developer.Name = strings.TrimSpace(document.Find(`.document-subtitle.primary`).Find(`[itemprop="name"]`).First().Text())
}

func (app *App) parsePrice(document *playStoreDocument) {
	content, ok := document.Find(`.details-info button.price.buy meta[itemprop="price"]`).Attr("content")
	if !ok {
		app.Price = "unknown"
		return
	}
	if content == "Install" {
		app.Price = "0"
	} else {
		app.Price = content
	}
}

func (app *App) parseOffersInApp(document *playStoreDocument) {
	app.OffersInApp = document.Find(`.inapp-msg`).Length() > 0
}

func (app *App) parseRating(document *playStoreDocument) {
	document.Find(`.rating-bar-container`).Each(func(i int, sel *goquery.Selection) {
		score := ParseInteger(strings.TrimSpace(sel.Find(`.bar-label`).First().Text()))
		rating := ParseInteger(strings.TrimSpace(sel.Find(`.bar-number`).First().Text()))
		app.Rating[strconv.Itoa(int(score))] = rating
	})
}

// Obtains the date of the last submitted apk.
func (app *App) parseUpdated(sel *goquery.Selection) {
	app.Updated = strings.TrimSpace(sel.Find(`.content`).First().Text())
}

func (app *App) parseSize(sel *goquery.Selection) {
	app.Size = strings.TrimSpace(sel.Find(`.content`).First().Text())
}

func (app *App) parseInstalls(sel *goquery.Selection) {
	app.Installs = strings.TrimSpace(sel.Find(`.content`).First().Text())
}

func (app *App) parseVersion(sel *goquery.Selection) {
	app.Version = strings.TrimSpace(sel.Find(`.content`).First().Text())
}

func (app *App) parseRequiresAndroid(sel *goquery.Selection) {
	app.RequiresAndroid = strings.TrimSpace(sel.Find(`.content`).First().Text())
}

func (app *App) parseContentRating(sel *goquery.Selection) {
	app.ContentRating = strings.TrimSpace(sel.Find(`.content`).First().Text())
}

func (app *App) parseContactDeveloper(sel *goquery.Selection) {
	sel.Find(`.content`).Find(`a`).Each(func(i int, anchor *goquery.Selection) {
		var ok bool
		switch strings.TrimSpace(anchor.Text()) {
		case "Visit Developer's Website":
			if app.Developer.Website, ok = anchor.Attr("href"); !ok {
				app.Developer.Website = ""
			}
		case "Email Developer":
			if app.Developer.Email, ok = anchor.Attr("href"); !ok {
				app.Developer.Email = ""
			}
			app.Developer.Email = strings.Replace(app.Developer.Email, "mailto:", "", 1)
		}
	})
}

// Checks if the current document corresponds to an application. In case there's
// no such thing it returns false.
func isValidApp(document *playStoreDocument) bool {
	return document.Find("title").First().Text() != "Not Found"
}
