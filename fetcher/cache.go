package fetcher

import "context"

type Cache struct {
	f Fetcher
	c map[string]string
}

func NewCache(f Fetcher) *Cache {
	return &Cache{
		f: f,
		c: map[string]string{},
	}
}

func (c *Cache) FetchRepo(ctx context.Context, module string) (string, error) {
	if v, ok := c.c[module]; ok {
		return v, nil
	}
	ret, err := c.f.FetchRepo(ctx, module)
	if err == nil {
		c.c[module] = ret
	}
	return ret, err
}
