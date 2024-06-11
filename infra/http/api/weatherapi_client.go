package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type weatherResponse struct {
	Location struct {
		Name           string  `json:"name"`
		Region         string  `json:"region"`
		Country        string  `json:"country"`
		Lat            float64 `json:"lat"`
		Lon            float64 `json:"lon"`
		TzID           string  `json:"tz_id"`
		LocaltimeEpoch int     `json:"localtime_epoch"`
		Localtime      string  `json:"localtime"`
	} `json:"location"`
	Current struct {
		LastUpdatedEpoch int     `json:"last_updated_epoch"`
		LastUpdated      string  `json:"last_updated"`
		TempC            float64 `json:"temp_c"`
		TempF            float64 `json:"temp_f"`
		IsDay            int     `json:"is_day"`
		Condition        struct {
			Text string `json:"text"`
			Icon string `json:"icon"`
			Code int    `json:"code"`
		} `json:"condition"`
		WindMph    float64 `json:"wind_mph"`
		WindKph    float64 `json:"wind_kph"`
		WindDegree int     `json:"wind_degree"`
		WindDir    string  `json:"wind_dir"`
		PressureMb float64 `json:"pressure_mb"`
		PressureIn float64 `json:"pressure_in"`
		PrecipMm   float64 `json:"precip_mm"`
		PrecipIn   float64 `json:"precip_in"`
		Humidity   int     `json:"humidity"`
		Cloud      int     `json:"cloud"`
		FeelslikeC float64 `json:"feelslike_c"`
		FeelslikeF float64 `json:"feelslike_f"`
		WindchillC float64 `json:"windchill_c"`
		WindchillF float64 `json:"windchill_f"`
		HeatindexC float64 `json:"heatindex_c"`
		HeatindexF float64 `json:"heatindex_f"`
		DewpointC  float64 `json:"dewpoint_c"`
		DewpointF  float64 `json:"dewpoint_f"`
		VisKm      float64 `json:"vis_km"`
		VisMiles   float64 `json:"vis_miles"`
		Uv         float64 `json:"uv"`
		GustMph    float64 `json:"gust_mph"`
		GustKph    float64 `json:"gust_kph"`
	} `json:"current"`
}

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type weatherClient struct {
	client HTTPClient
	apiKey string
}

func NewWeatherClient(client HTTPClient, apiKey string) (*weatherClient, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("error config client")
	}

	return &weatherClient{
		client: client,
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

	resp, err := w.client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		return GetTempOutput{}, fmt.Errorf("fail to do request: %w", err)
	}

	respPayload := &weatherResponse{}
	respBuf := bytes.Buffer{}
	_, err = respBuf.ReadFrom(resp.Body)
	if err != nil {
		return GetTempOutput{}, fmt.Errorf("fail to read response: %w", err)
	}

	err = json.Unmarshal(respBuf.Bytes(), respPayload)
	if err != nil {
		return GetTempOutput{}, fmt.Errorf("fail to decode response: %w, full response: %s", err, respBuf.String())
	}

	return GetTempOutput{
		TempC: respPayload.Current.TempC,
		TempF: respPayload.Current.TempF,
	}, nil
}
