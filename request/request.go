package request

import (
	"context"
	"net/http"
)

func DoNewRequestWithContext(ctx context.Context, url string) (*http.Response, error) {
	res, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := http.DefaultClient.Do(res)
	if err != nil {
		return nil, err
	}
	return resp, err
}
