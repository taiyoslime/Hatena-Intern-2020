package fetcher

import (
	"context"
)

// Fetch は受け取った文書を HTML に変換する
func Fetch(ctx context.Context, url string) (string, error) {
	return url, nil 
}
