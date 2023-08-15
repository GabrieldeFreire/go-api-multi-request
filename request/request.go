package request

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/GabrieldeFreire/multithreading/log"
)

var logger *slog.Logger = log.GetInstance()

type ChanResp struct {
	Url  string
	Resp *http.Response
	Err  error
}

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
