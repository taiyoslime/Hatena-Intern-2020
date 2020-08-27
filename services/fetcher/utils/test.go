package utils

type CacheClient interface {
	Get(string) (interface{}, bool)
	SetDefault(string, interface{})
}

type TestCacheClient struct {
	CacheClient
	table    map[string]string
	CacheHit int
}

func CreateTestCacheClient() *TestCacheClient {
	return &TestCacheClient{
		table:    map[string]string{},
		CacheHit: 0,
	}
}

func (c *TestCacheClient) Get(key string) (interface{}, bool) {
	val, ok := c.table[key]
	if ok {
		c.CacheHit++
	}
	return val, ok
}

func (c *TestCacheClient) SetDefault(key string, val interface{}) {
	c.table[key] = val.(string)
}
