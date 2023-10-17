package schema

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/GabrieldeFreire/multithreading/request"
)

const (
	BrasilApiUrl = "https://brasilapi.com.br/api/cep/v1/%s"
)

type ApiBrasilApi struct {
	Cep        string `json:"cep"`
	Uf         string `json:"state"`
	Localidade string `json:"city"`
	Bairro     string `json:"neighborhood"`
	Logradouro string `json:"street"`
}

type BrasilApi struct {
	Api              *ApiBrasilApi
	Cep              string
	ctx              context.Context
	cancel           context.CancelFunc
	reqChan          chan *ApiCepInfo
	statusOK         bool
	name             string
	DeadlineExceeded bool
}

func (b BrasilApi) Name() string {
	return b.name
}

func (b BrasilApi) StatusOk() bool {
	return b.statusOK
}

func (b *BrasilApi) CancelContext() {
	b.cancel()
}

func (b *BrasilApi) ToApiCep() *ApiCepInfo {
	return NewApiCep(
		b.Api.Cep,
		b.Api.Uf,
		b.Api.Localidade,
		b.Api.Bairro,
		b.Api.Logradouro,
		b.name,
		b.statusOK,
		b.DeadlineExceeded,
	)
}

func (b *BrasilApi) DoRequest() {
	urlCep := fmt.Sprintf(BrasilApiUrl, b.Cep)
	resp, err := request.DoNewRequestWithContext(b.ctx, urlCep)

	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			b.DeadlineExceeded = true
			b.reqChan <- b.ToApiCep()
			return
		}
		b.reqChan <- b.ToApiCep()
	} else if resp.StatusCode != http.StatusOK {
		b.reqChan <- b.ToApiCep()
	} else {
		err = json.NewDecoder(resp.Body).Decode(&b.Api)
		if err != nil {
			b.reqChan <- b.ToApiCep()
			return
		}
		b.statusOK = true
		b.reqChan <- b.ToApiCep()
	}
}

func NewBrasilApi(cep string, reqChan chan *ApiCepInfo) CepInterface {
	ctx, cancel := context.WithTimeout(context.Background(), REQUEST_MAX_DURATION)
	return &BrasilApi{
		Api:              &ApiBrasilApi{},
		Cep:              formatCepWithDash(cep),
		ctx:              ctx,
		cancel:           cancel,
		reqChan:          reqChan,
		statusOK:         false,
		name:             "BrasilApi",
		DeadlineExceeded: false,
	}
}
