package api

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_weatherClient_GetTemp_Success(t *testing.T) {

	w := &weatherClient{}
	_, err := w.GetTemp(context.TODO(), "Balneario Gaivota, SC, Brasil")

	assert.NoError(t, err)
}

func Test_weatherClient_GetTemp(t *testing.T) {

	w := &weatherClient{}
	got, err := w.GetTemp(context.TODO(), "Balneario Gaivota, SC, Brasil")
	if err != nil {
		t.Failed()
	}
	assert.NotEqual(t, float64(0), got.TempC)
}
