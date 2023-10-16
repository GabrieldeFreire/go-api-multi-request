package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/GabrieldeFreire/multithreading/log"
	"github.com/GabrieldeFreire/multithreading/schema"
)

const (
	REQUEST_MAX_DURATION = 1 * time.Second
)

var CEP_APIS = [...]func(cep string, reqChan chan *schema.ApiCepInfo) schema.CepInterface{
	schema.NewBrasilApi,
	schema.NewViaCep,
}

var logger *slog.Logger = log.GetInstance()

func main() {
	var cep string
	flag.StringVar(&cep, "cep", "09530-210", "Cep para busca")
	flag.Parse()

	getCep(cep)
}

func getCep(cep string) {
	respChan := make(chan *schema.ApiCepInfo)

	allApiStructs := map[string]schema.CepInterface{}
	mutex := &sync.RWMutex{}

	for _, apiFunc := range CEP_APIS {
		apiFunc := apiFunc
		go func() {
			apiFunc := apiFunc
			apiStruct := apiFunc(cep, respChan)
			mutex.Lock()
			allApiStructs[apiStruct.Name()] = apiStruct
			mutex.Unlock()
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
