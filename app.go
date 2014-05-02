package playstore

import (
	"github.com/PuerkitoBio/goquery"
	"strconv"
	"strings"
)

type App struct {
	Title                string            `json:"title"`
	Author               string            `json:"author"`
	Category             string            `json:"category"`
	OffersInApp          bool              `json:"offersInApp"`
	AverageRating        float64           `json:"averageRating"`
	Rating               map[string]int64  `json:"rating"`
	Reviews              int64             `json:"reviews"`
	ScreenshotUrls       []string          `json:"screenshotUrls"`
	TabletScreenshotUrls []string          `json:"tabletScreenshotUrls"`
	Updated              string            `json:"updated"`
	Version              string            `json:"version"`
	Price                float64           `json:"price"`
	Size                 string            `json:"size"`
	RequiresAndroid      string            `json:"requiresAndroid"`
	ContentRating        string            `json:"contentRating"`
	Installs             string            `json:"installs"`
	PlainDescription     map[string]string `json:"plainDescription"`
	HtmlDescription      map[string]string `json:"htmlDescription"`
	Developer            struct {
		Email   string `json:"email"`
		Website string `json:"website"`
	} `json:"developer"`
}

func NewApp() *App {
	app := &App{}
	app.Rating = make(map[string]int64, 5)
	app.PlainDescription = map[string]string{}
	app.HtmlDescription = map[string]string{}
	return app
}

// Looks through the document to extract the app info.
func parseApp(document *playStoreDocument, lang string) (*App, error) {
	app := NewApp()
	app.parseTitle(document)
	app.parseAverageRating(document)
	app.parseReviews(document)
	app.parseMetaInfo(document)
	app.parseDescription(document, lang)
	app.parseScreenshotUrls(document)
	app.parseCategory(document)
	app.parseRating(document)
	return app, nil
}

func (app *App) parseTitle(document *playStoreDocument) {
	title := document.Find(`.document-title`).Children().First().Text()
	app.Title = strings.TrimSpace(title)
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

func (app *App) parseCategory(document *playStoreDocument) {
	app.Category = strings.TrimSpace(document.Find(`.category`).Find(`[itemprop="genre"]`).First().Text())
}

func (app *App) parseRating(document *playStoreDocument) {
	document.Find(`.rating-bar-container`).Each(func(i int, sel *goquery.Selection) {
		score := ParseInteger(strings.TrimSpace(sel.Find(`.bar-label`).First().Text()))
		rating := ParseInteger(strings.TrimSpace(sel.Find(`.bar-number`).First().Text()))
		app.Rating[string(score)] = rating
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

func isValidApp(document *playStoreDocument) bool {
	return document.Find("title").First().Text() != "Not Found"
}
