package fetcher

import (
	"context"
	utils "fetcher/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

type FetchTestCase struct {
	in       string
	out      string
	err      bool
	cacheHit int // キャッシュにヒットする数
}

var fetchTestCases = []FetchTestCase{
	{
		in:  "https://google.com",
		out: "Google",
		err: false,
	},
	{
		// <title>が無い場合
		in:  "http://dev.taiyosli.me",
		out: "",
		err: false,
	},
	{
		// 存在しないurlの場合
		in:  "https://url.which.does.not.exist",
		out: "",
		err: true,
	},
	/*
		http://dev.taiyosli.me/robots.txt

		User-agent: *
		Disallow: /disallow
		Allow: /disallow/allow
	*/
	{
		in:  "http://dev.taiyosli.me/disallow",
		out: "",
		err: true,
	},
	{
		in:  "http://dev.taiyosli.me/disallow/allow",
		out: "OK",
		err: false,
	},
	{
		in:       "http://dev.taiyosli.me/disallow/allow",
		out:      "OK",
		err:      false,
		cacheHit: 1,
	},
}

func Test_Fetch(t *testing.T) {
	testCacheClient := utils.CreateTestCacheClient()

	for _, testCase := range fetchTestCases {
		url := testCase.in
		title, err := Fetch(context.Background(), testCacheClient, url)
		if !testCase.err {
			assert.NoError(t, err)
		} else {
			assert.Error(t, err)
		}
		if testCase.cacheHit != 0 {
			assert.Equal(t, testCacheClient.CacheHit, testCase.cacheHit)
		}
		assert.Equal(t, title, testCase.out)
	}
}
