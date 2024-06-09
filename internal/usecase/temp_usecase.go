package usecase

import (
	"context"
	"fmt"
	"github.com/airtongit/fc-temp-by-cep/infra/http/api"
)

type tempUsecase struct {
	tempClient TempClient
}

func NewTempUsecase(tempClient TempClient) *tempUsecase {
	return &tempUsecase{
		tempClient: tempClient,
	}
}

type TempClient interface {
	GetTemp(ctx context.Context, q string) (api.GetTempOutput, error)
}

type TempUsecaseInput struct {
	Localidade string
	Uf         string
	Pais       string
}

type TempUsecaseOutput struct {
	TempC float64
	TempF float64
}

func (t *tempUsecase) Execute(ctx context.Context, input TempUsecaseInput) (TempUsecaseOutput, error) {

	q := fmt.Sprintf("%s, %s, %s", input.Localidade, input.Uf, input.Pais)
	tempOutput, err := t.tempClient.GetTemp(ctx, q)
	if err != nil {
		return TempUsecaseOutput{}, err
	}
	return TempUsecaseOutput{
		TempC: tempOutput.TempC,
		TempF: tempOutput.TempF,
	}, nil
}
