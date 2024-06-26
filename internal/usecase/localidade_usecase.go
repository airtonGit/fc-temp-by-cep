package usecase

import (
	"context"
	"errors"
	"fmt"
	"github.com/airtongit/fc-temp-by-cep/infra/http/api"
)

type CepClient interface {
	GetLocalidade(ctx context.Context, cep string) (*api.LocalidadeResponse, error)
}

var (
	ErrCepNotFound = fmt.Errorf("cep not found")
)

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
		if errors.Is(err, ErrCepNotFound) {
			return LocalidadeOutput{}, ErrCepNotFound
		}
		return LocalidadeOutput{}, err
	}

	if localidadeOutput.Erro {
		return LocalidadeOutput{}, fmt.Errorf("can not find zipcode")
	}

	return LocalidadeOutput{
		Localidade: localidadeOutput.Localidade,
		Uf:         localidadeOutput.Uf,
		Pais:       "Brasil",
	}, nil
}
