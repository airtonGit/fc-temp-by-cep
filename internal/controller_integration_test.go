package internal

import (
	"context"
	"net/http"
	"os"
	"testing"

	"github.com/airtongit/fc-temp-by-cep/infra/http/api"
	"github.com/airtongit/fc-temp-by-cep/internal/usecase"
)

func TestController(t *testing.T) {
	if os.Getenv("WEATHER") == "" {
		t.Skip("skipping integration test, WEATHER environment variable not set")
	}

	cepClient := api.NewViaCEPClient("http://viacep.com.br")
	localidadeUsecase := usecase.NewLocalidadeUsecase(cepClient)
	tempClient, _ := api.NewWeatherClient(&http.Client{}, os.Getenv("WEATHER"))
	tempUsecase := usecase.NewTempUsecase(tempClient)
	kelvinService := usecase.NewKelvinService()
	tempByCEPctrl := NewTempByLocaleController(localidadeUsecase, tempUsecase, kelvinService)

	temp, err := tempByCEPctrl.GetTemp(context.Background(), "88955-000")
	if err != nil {
		t.Errorf("GetTemp() error = %v", err)
	}
	t.Error(temp)
}
