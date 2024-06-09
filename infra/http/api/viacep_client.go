package api

import (
	"bytes"
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
	Erro        bool   `json:"erro"`
}

func (v *viaCEPClient) GetLocalidade(ctx context.Context, cep string) (*LocalidadeResponse, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/ws/%s/json", v.baseURL, cep), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	respBuf := bytes.Buffer{}
	respBuf.ReadFrom(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}
	cepPayload := new(LocalidadeResponse)
	err = json.Unmarshal(respBuf.Bytes(), cepPayload)
	if err != nil {
		return nil, fmt.Errorf("failed to decode response: %w, full response: %s", err, respBuf.String())
	}

	return cepPayload, nil
}
