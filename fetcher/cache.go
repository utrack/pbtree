package fetcher

import "context"

type Cache struct {
	f Fetcher
	c map[string]FileOpener
}

func NewCache(f Fetcher) *Cache {
	return &Cache{
		f: f,
		c: map[string]FileOpener{},
	}
}

func (c *Cache) FetchRepo(ctx context.Context, module string) (FileOpener, error) {
	if v, ok := c.c[module]; ok {
		return v, nil
	}
	ret, err := c.f.FetchRepo(ctx, module)
	if err == nil {
		c.c[module] = ret
	}
	return ret, err
}
