package usecase

import (
	"context"
	"github.com/airtongit/fc-temp-by-cep/infra/http/api"
)

type CepClient interface {
	GetLocalidade(ctx context.Context, cep string) (*api.LocalidadeResponse, error)
}

type localidadeUsecase struct {
	cepClient CepClient
}

type LocalidadeInput struct {
	Cep string
}

type LocalidadeOutput struct {
	Localidade string
	Uf         string
	Pais       string
}

func NewLocalidadeUsecase(cepClient CepClient) *localidadeUsecase {
	return &localidadeUsecase{
		cepClient: cepClient,
	}
}

func (l *localidadeUsecase) Execute(ctx context.Context, input LocalidadeInput) (LocalidadeOutput, error) {
	localidadeOutput, err := l.cepClient.GetLocalidade(ctx, input.Cep)
	if err != nil {
		return LocalidadeOutput{}, err
	}

	return LocalidadeOutput{
		Localidade: localidadeOutput.Localidade,
		Uf:         localidadeOutput.Uf,
		Pais:       "Brasil",
	}, nil
}
