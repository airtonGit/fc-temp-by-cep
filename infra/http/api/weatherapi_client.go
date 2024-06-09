package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

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

func NewWeatherClient(apiKey string) (*weatherClient, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("error config client")
	}

	return &weatherClient{
		apiKey: apiKey,
	}, nil
}

type GetTempOutput struct {
	TempC float64 `json:"temp_c"`
	TempF float64 `json:"temp_f"`
}

func (w *weatherClient) GetTemp(ctx context.Context, q string) (GetTempOutput, error) {
	urlBase := fmt.Sprintf("http://api.weatherapi.com/v1/current.json?key=%s&q=%s&aqi=no", w.apiKey, url.QueryEscape(q))
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, urlBase, nil)
	if err != nil {
		return GetTempOutput{}, fmt.Errorf("fail to create request: %w", err)
	}

	client := &http.Client{}

	resp, err := client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		return GetTempOutput{}, fmt.Errorf("fail to do request: %w", err)
	}

	respPayload := &weatherResponse{}

	respBuf := make([]byte, 1024)
	_, err = resp.Body.Read(respBuf)
	if err != nil {
		return GetTempOutput{}, fmt.Errorf("fail to load response: %w", err)
	}

	err = json.Unmarshal(respBuf, respPayload)
	if err != nil {
		return GetTempOutput{}, fmt.Errorf("fail to decode response: %w, full response: %s", err, string(respBuf))
	}

	return GetTempOutput{
		TempC: respPayload.Current.TempC,
		TempF: respPayload.Current.TempF,
	}, nil
}
