package fetcher

import (
	"context"
	"errors"
	"net/http"
	_url "net/url"

	"github.com/PuerkitoBio/goquery"
	"github.com/temoto/robotstxt"
)

type CacheClient interface {
	Get(string) (interface{}, bool)
	SetDefault(string, interface{})
}

// Fetch は受け取った文書を HTML に変換する
func Fetch(ctx context.Context, cache CacheClient, url string) (string, error) {
	client := &http.Client{}

	// キャッシュにあるならそれを返すようにする
	cachedTitle, ok := cache.Get(url)
	if ok {
		return cachedTitle.(string), nil
	}

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	// robots.txtの内容を読み，disallowであればリクエストを行わない
	if !isAllowed(ctx, url, client) {
		return "", errors.New("the site's robot.txt denies request")
	}

	res, err := client.Do(request)
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
	cache.SetDefault(url, title)
	return title, nil
}

func isAllowed(ctx context.Context, url string, client *http.Client) bool {

	u, err := _url.Parse(url)
	if err != nil {
		return false
	}
	host := u.Host
	scheme := u.Scheme
	path := u.Path

	req, err := http.NewRequest("GET", scheme+"://"+host+"/robots.txt", nil)
	if err != nil {
		return false
	}

	res, err := client.Do(req)
	if err != nil {
		return false
	}
	defer res.Body.Close()

	if res.StatusCode == 200 {
		robots, err := robotstxt.FromResponse(res)
		if err == nil {
			group := robots.FindGroup(req.UserAgent())
			allow := group.Test(path)
			if !allow {
				return false
			}
		}
	}
	return true
}
