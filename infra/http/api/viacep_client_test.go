package api

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_viaCEPClient_Ask_Success(t *testing.T) {

	// given
	v := &viaCEPClient{
		baseURL: "http://viacep.com.br",
	}

	// When
	_, err := v.GetLocalidade(context.TODO(), "88955-000")

	// Then
	assert.NoError(t, err)
}

func Test_viaCEPClient_Ask_Response(t *testing.T) {

	// given
	v := &viaCEPClient{
		baseURL: "http://viacep.com.br",
	}

	// When
	got, err := v.GetLocalidade(context.TODO(), "88955-000")
	assert.NoError(t, err)

	// Then
	assert.Equal(t, "SC", got.Uf, "Esperado SC, got=%s", got.Uf)
}

func Test_viaCEPClient_CEP_NotFound(t *testing.T) {

	// given
	v := &viaCEPClient{
		baseURL: "http://viacep.com.br",
	}

	// When
	got, err := v.GetLocalidade(context.TODO(), "12345678")
	assert.NoError(t, err)

	// Then
	assert.Equal(t, true, got.Erro, "Error is true, cep not found")
}
