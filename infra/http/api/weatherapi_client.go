package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

const baseUrl = "https://www.weatherapi.com/docs/#"

const weatherAPIKey = "955781466c1e414e9e9181300240806"

type weatherResponse struct {
	Current struct {
		LastUpdatedEpoch int     `json:"last_updated_epoch"`
		LastUpdated      string  `json:"last_updated"`
		TempC            float64 `json:"temp_c"`
		TempF            float64 `json:"temp_f"`
		IsDay            int     `json:"is_day"`
	} `json:"current"`
}

type weatherClient struct {
	apiKey string
}

func NewWeatherClient(apiKey string) *weatherClient {
	return &weatherClient{
		apiKey: apiKey,
	}
}

type GetTempOutput struct {
	TempC float64 `json:"temp_c"`
	TempF float64 `json:"temp_f"`
}

func (w *weatherClient) GetTemp(ctx context.Context, q string) (GetTempOutput, error) {
	urlBase := fmt.Sprintf("http://api.weatherapi.com/v1/current.json?key=%s&q=%s&aqi=no", w.apiKey, url.QueryEscape(q))
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, urlBase, nil)
	if err != nil {
		return GetTempOutput{}, err
	}

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return GetTempOutput{}, err
	}

	respPayload := &weatherResponse{}
	err = json.NewDecoder(resp.Body).Decode(respPayload)
	if err != nil {
		return GetTempOutput{}, err
	}
	return GetTempOutput{
		TempC: respPayload.Current.TempC,
		TempF: respPayload.Current.TempF,
	}, nil
}
