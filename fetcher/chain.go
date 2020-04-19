package fetcher

import "context"

type Chainer struct {
	ff []Fetcher
}

func Chain(ff ...Fetcher) Chainer {
	return Chainer{ff: ff}
}

func (c Chainer) FetchRepo(ctx context.Context, module string) (string, error) {
	var ret string
	var err error
	for _, f := range c.ff {
		ret, err = f.FetchRepo(ctx, module)
		if err == ErrOtherFetcher {
			continue
		}
		return ret, err
	}
	return ret, err
}
