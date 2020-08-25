package fetcher

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

type FetchTestCase struct {
	in  string
	out string
	err bool
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
}

func Test_Fetch(t *testing.T) {
	for _, testCase := range fetchTestCases {
		url := testCase.in
		title, err := Fetch(context.Background(), url)
		if !testCase.err {
			assert.NoError(t, err)
		} else {
			assert.Error(t, err)
		}

		assert.Equal(t, title, testCase.out)
	}
}
