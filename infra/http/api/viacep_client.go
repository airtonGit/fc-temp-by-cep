package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type viaCEPClient struct {
	baseURL string
}

func NewViaCEPClient(baseURL string) *viaCEPClient {
	return &viaCEPClient{
		baseURL: baseURL,
	}
}

type LocalidadeResponse struct {
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

func (v *viaCEPClient) GetLocalidade(ctx context.Context, cep string) (*LocalidadeResponse, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/ws/%s/json", v.baseURL, cep), nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	cepPayload := new(LocalidadeResponse)
	err = json.NewDecoder(resp.Body).Decode(cepPayload)
	if err != nil {
		return nil, err
	}

	return cepPayload, nil
}
