package internal

import (
	"context"
	"github.com/airtongit/fc-temp-by-cep/internal/usecase"
)

type LocalidadeUsecase interface {
	Execute(ctx context.Context, input usecase.LocalidadeInput) (usecase.LocalidadeOutput, error)
}

type TempUsecase interface {
	Execute(ctx context.Context, input usecase.TempUsecaseInput) (usecase.TempUsecaseOutput, error)
}

type KelvinService interface {
	GetKelvin(tempC float64) float64
}

type tempByLocaleController struct {
	localidadeUsecase LocalidadeUsecase
	tempUsecase       TempUsecase
	kelvinService     KelvinService
}

func NewTempByLocaleController(localidadeUsecase LocalidadeUsecase, tempUsecase TempUsecase, kelvinService KelvinService) *tempByLocaleController {
	return &tempByLocaleController{
		localidadeUsecase: localidadeUsecase,
		tempUsecase:       tempUsecase,
		kelvinService:     kelvinService,
	}
}

type Temp struct {
	TempC float64 `json:"temp_C,omitempty"`
	TempF float64 `json:"temp_F,omitempty"`
	TempK float64 `json:"temp_K,omitempty"`
}

func (t *tempByLocaleController) GetTemp(ctx context.Context, cep string) (Temp, error) {

	localidadeInput := usecase.LocalidadeInput{
		Cep: cep,
	}
	localidade, err := t.localidadeUsecase.Execute(ctx, localidadeInput)
	if err != nil {
		return Temp{}, err
	}

	tempUsecaseInput := usecase.TempUsecaseInput{
		Localidade: localidade.Localidade,
		Uf:         localidade.Uf,
		Pais:       localidade.Pais,
	}
	temp, err := t.tempUsecase.Execute(ctx, tempUsecaseInput)
	if err != nil {
		return Temp{}, err
	}

	kelvin := t.kelvinService.GetKelvin(temp.TempC)

	return Temp{
		TempC: temp.TempC,
		TempF: temp.TempF,
		TempK: kelvin,
	}, nil
}
