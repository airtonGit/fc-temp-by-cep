package http

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

const baseURL = "http://viacep.com.br"

type viaCEPClient struct {
	baseURL string
}

type CepResponse struct {
	Cep         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	Uf          string `json:"uf"`
	Unidade     string `json:"unidade"`
	Ibge        string `json:"ibge"`
	Gia         string `json:"gia"`
}

func (v *viaCEPClient) Ask(ctx context.Context, cep string) (*CepResponse, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/ws/%s/json", v.baseURL, cep), nil)

	if err != nil {
		return nil, err
	}

	var cepPayload *CepResponse
	err = json.NewDecoder(req.Body).Decode(cepPayload)
	if err != nil {
		return nil, err
	}

	return cepPayload, nil
}
