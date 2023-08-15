package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log/slog"
	"time"

	"github.com/GabrieldeFreire/multithreading/log"
	"github.com/GabrieldeFreire/multithreading/schema"
)

const (
	REQUEST_MAX_DURATION = 10 * time.Second
)

var CEP_APIS = [...]func(cep string, reqChan chan *schema.ApiCepInfo) schema.CepInterface{
	schema.NewBrasilApi,
	schema.NewViaCep,
}

var logger *slog.Logger = log.GetInstance()

// "https://cdn.apicep.com/file/apicep/%s.json",
// "http://viacep.com.br/ws/%s/json/",

func main() {
	// cep = flag.String("cep", "09530-210", "current environment")

	var cep string
	flag.StringVar(&cep, "cep", "09530-210", "Cep para busca")
	flag.Parse()
	getCep(cep)
}

func getCep(cep string) {
	respChan := make(chan *schema.ApiCepInfo)

	allApiStructs := map[string]schema.CepInterface{}

	for _, apiFunc := range CEP_APIS {
		apiFunc := apiFunc
		go func() {
			apiFunc := apiFunc
			apiStruct := apiFunc(cep, respChan)
			allApiStructs[apiStruct.Name()] = apiStruct
			apiStruct.DoRequest()
		}()
	}

	count := 0
	var apiCepResponse *schema.ApiCepInfo
urlLoop:
	for {
		select {
		case apiCepResponse = <-respChan:
			count++
			if !apiCepResponse.StatusOK {
				logger.Error(fmt.Sprintf("%s request failed", apiCepResponse.ApiName))
				continue
			}
			close(respChan)
			for name, apiStruct := range allApiStructs {
				if name == apiCepResponse.ApiName {
					break
				}
				apiStruct.CancelContext()
			}
			break urlLoop
		default:
			if count == len(CEP_APIS) {
				fmt.Println("all requests failed")
				return
			}

		}
	}

	prettyPrint(apiCepResponse.ApiName, apiCepResponse.Api)
}

func prettyPrint(apiName string, apiStruct *schema.ApiCep) {
	empJSON, err := json.MarshalIndent(apiStruct, "", "  ")
	if err != nil {
		logger.Error(err.Error())
	}
	fmt.Printf("%s: %s\n", apiName, string(empJSON))
}
