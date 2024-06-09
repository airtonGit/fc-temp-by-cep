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
	Location struct {
		Name           string  `json:"name"`
		Region         string  `json:"region"`
		Country        string  `json:"country"`
		Lat            float64 `json:"lat"`
		Lon            float64 `json:"lon"`
		TzId           string  `json:"tz_id"`
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
	urlBase := fmt.Sprintf("http://api.weatherapi.com/v1/current.json?key=%s&q=%s&aqi=no", weatherAPIKey, url.QueryEscape(q))
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
