package api

import (
	"context"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_weatherClient_GetTemp(t *testing.T) {

	if os.Getenv("WEATHER") == "" {
		t.Skip("skipping integration test, need env var WEATHER")
	}

	w, err := NewWeatherClient(&http.Client{}, os.Getenv("WEATHER"))
	if err != nil {
		t.Errorf("test need env var")
	}
	got, err := w.GetTemp(context.TODO(), "Balneario Gaivota, SC, Brasil")
	if err != nil {
		t.Failed()
	}
	assert.NotEqual(t, float64(0), got.TempC)
}
