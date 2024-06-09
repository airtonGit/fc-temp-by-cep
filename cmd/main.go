package main

import (
	"fmt"
	"github.com/airtongit/fc-temp-by-cep/infra/http/api"
	"github.com/airtongit/fc-temp-by-cep/infra/http/handler"
	"github.com/airtongit/fc-temp-by-cep/internal/usecase"
	"net/http"
	"os"

	"github.com/airtongit/fc-temp-by-cep/internal"
	"github.com/go-chi/chi"
	"github.com/joho/godotenv"
)

func main() {

	if _, err := os.Stat(".env"); err == nil {
		err := godotenv.Load()
		if err != nil {
			fmt.Println("Error loading .env file", err)
			return
		}
	}

	cepClient := api.NewViaCEPClient("http://viacep.com.br")
	localidadeUsecase := usecase.NewLocalidadeUsecase(cepClient)
	tempClient := api.NewWeatherClient(os.Getenv("WEATHER"))
	tempUsecase := usecase.NewTempUsecase(tempClient)
	kelvinService := usecase.NewKelvinService()
	tempByCEPctrl := internal.NewTempByLocaleController(localidadeUsecase, tempUsecase, kelvinService)

	r := chi.NewRouter()
	r.Get("/{cep}", handler.CepHandler(tempByCEPctrl))

	fmt.Println("Listening on :8080")
	err := http.ListenAndServe(":8080", r)
	if err != nil {
		fmt.Println(err)
		return
	}
}
