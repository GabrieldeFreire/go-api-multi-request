package schema

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/GabrieldeFreire/multithreading/request"
)

const (
	ViaCepURL = "http://viacep.com.br/ws/%s/json/"
)

type ApiViaCep struct {
	Cep        string `json:"cep"`
	Uf         string `json:"uf"`
	Localidade string `json:"localidade"`
	Bairro     string `json:"bairro"`
	Logradouro string `json:"logradouro"`
}

type ViaCep struct {
	Api      *ApiViaCep
	Cep      string
	ctx      context.Context
	cancel   context.CancelFunc
	reqChan  chan *ApiCepInfo
	statusOK bool
	name     string
}

func (v ViaCep) Name() string {
	return v.name
}

func (v *ViaCep) CancelContext() {
	v.cancel()
}

func (v *ViaCep) ToApiCep() *ApiCepInfo {
	return NewApiCep(
		v.Api.Cep,
		v.Api.Uf,
		v.Api.Localidade,
		v.Api.Bairro,
		v.Api.Logradouro,
		v.name,
		v.statusOK,
	)
}

func (v *ViaCep) DoRequest() {
	urlCep := fmt.Sprintf(ViaCepURL, v.Cep)
	resp, err := request.DoNewRequestWithContext(v.ctx, urlCep)

	if err != nil {
		v.reqChan <- v.ToApiCep()
	} else if resp.StatusCode != http.StatusOK {
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

func NewViaCep(cep string, reqChan chan *ApiCepInfo) CepInterface {
	ctx, cancel := context.WithTimeout(context.Background(), REQUEST_MAX_DURATION)
	return &ViaCep{
		Api:      &ApiViaCep{},
		Cep:      formatCepWithDash(cep),
		ctx:      ctx,
		cancel:   cancel,
		reqChan:  reqChan,
		statusOK: false,
		name:     "ViaCep",
	}
}
