package fetcher

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Fetch(t *testing.T) {
	url := "https://google.com"
	html, err := Fetch(context.Background(), url)
	assert.NoError(t, err)
	assert.Equal(t, `Google`, html)
}
