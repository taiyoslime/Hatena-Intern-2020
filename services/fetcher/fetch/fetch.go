package fetcher

import (
	"context"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)


// Fetch は受け取った文書を HTML に変換する
func Fetch(ctx context.Context, url string) (string, error) {
	res, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return "", err
	}
	
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return "", err
	}
	
	title := doc.Find("title").Text()
	return title, nil 
}
