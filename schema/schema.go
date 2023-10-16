package schema

import (
	"time"
)

const (
	REQUEST_MAX_DURATION = 10 * time.Second
)

type CepInterface interface {
	ToApiCep() *ApiCepInfo
	Name() string
	DoRequest()
	CancelContext()
}

type ApiCep struct {
	Cep        string
	Uf         string
	Localidade string
	Bairro     string
	Logradouro string
}

type ApiCepInfo struct {
	Api      *ApiCep
	ApiName  string
	StatusOK bool
}

func NewApiCep(cep, uf, localidade, bairro, logradouro, apiName string, statusOk bool) *ApiCepInfo {
	return &ApiCepInfo{
		&ApiCep{
			Cep:        formatCepWithDash(cep),
			Uf:         uf,
			Localidade: localidade,
			Bairro:     bairro,
			Logradouro: logradouro,
		},
		apiName,
		statusOk,
	}
}

func formatCepWithDash(cep string) string {
	if cep == "" {
		return cep
	}
	if len(cep) == 9 {
		return cep
	}
	return cep[:5] + "-" + cep[5:]
}
