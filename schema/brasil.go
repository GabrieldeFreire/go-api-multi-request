package schema

import (
	"context"
	"encoding/json"
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
	Api      *ApiBrasilApi
	Cep      string
	ctx      context.Context
	cancel   context.CancelFunc
	reqChan  chan *ApiCepInfo
	statusOK bool
	name     string
}

func (v BrasilApi) Name() string {
	return v.name
}

func (v BrasilApi) StatusOk() bool {
	return v.statusOK
}

func (v *BrasilApi) CancelContext() {
	v.cancel()
}

func (e *BrasilApi) ToApiCep() *ApiCepInfo {
	return NewApiCep(
		e.Api.Cep,
		e.Api.Uf,
		e.Api.Localidade,
		e.Api.Bairro,
		e.Api.Logradouro,
		e.name,
		e.statusOK,
	)
}

func (v *BrasilApi) DoRequest() {
	urlCep := fmt.Sprintf(BrasilApiUrl, v.Cep)
	resp, err := request.DoNewRequestWithContext(v.ctx, urlCep)

	if err != nil {
		v.reqChan <- v.ToApiCep()
	} else if (resp.StatusCode < http.StatusEarlyHints) || (resp.StatusCode > http.StatusMultipleChoices) {
		v.reqChan <- v.ToApiCep()
	} else {
		err = json.NewDecoder(resp.Body).Decode(&v.Api)
		if err != nil {
			v.reqChan <- v.ToApiCep()
			return
		}
		v.statusOK = true
		v.reqChan <- v.ToApiCep()
	}
}

func NewBrasilApi(cep string, reqChan chan *ApiCepInfo) CepInterface {
	ctx, cancel := context.WithTimeout(context.Background(), REQUEST_MAX_DURATION)
	return &BrasilApi{
		Api:      &ApiBrasilApi{},
		Cep:      formatCepWithDash(cep),
		ctx:      ctx,
		cancel:   cancel,
		reqChan:  reqChan,
		statusOK: false,
		name:     "BrasilApi",
	}
}
