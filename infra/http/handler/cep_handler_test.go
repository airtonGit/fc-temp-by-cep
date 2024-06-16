// handlers_test.go
package handler

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/airtongit/fc-temp-by-cep/infra/http/api"
	"github.com/airtongit/fc-temp-by-cep/internal"
	"github.com/airtongit/fc-temp-by-cep/internal/usecase"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/mock"
	"go.opentelemetry.io/otel/trace"
)

type traceSpan struct {
	mock.Mock
}

func (t *traceSpan) End(options ...trace.SpanEndOption) {
	t.Called()
}

type tracerMock struct {
	mock.Mock
	span *traceSpan
}

func (t *tracerMock) Start(ctx context.Context, spanName string, opts ...trace.SpanStartOption) (context.Context, internal.TraceSpan) {
	t.Called(ctx, spanName, opts)
	return context.TODO(), t.span
}

func TestHealthCepHandlerInvalidCEP(t *testing.T) {

	cepClient := api.NewViaCEPClient("http://viacep.com.br")
	localidadeUsecase := usecase.NewLocalidadeUsecase(cepClient)

	os.Setenv("WEATHER", "955781466c1e414e9e9181300240806")

	tempClient, err := api.NewWeatherClient(&http.Client{}, os.Getenv("WEATHER"))
	if err != nil {
		t.Fatal(err)
		return
	}

	myTracerMock := new(tracerMock)
	myTracerMock.On("Start", mock.Anything, mock.AnythingOfType("string")).Return()

	tempUsecase := usecase.NewTempUsecase(tempClient)
	kelvinService := usecase.NewKelvinService()
	tempByCEPctrl := internal.NewTempByLocaleController(myTracerMock, localidadeUsecase, tempUsecase, kelvinService)

	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	req, err := http.NewRequest("GET", "/cep/011530000", nil)
	if err != nil {
		t.Fatal(err)
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := MakeCepHandler(tempByCEPctrl)
	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)
	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusUnprocessableEntity {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected := `invalid zipcode`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestHealthCepHandlerValidCEP(t *testing.T) {

	cepClient := api.NewViaCEPClient("http://viacep.com.br")
	localidadeUsecase := usecase.NewLocalidadeUsecase(cepClient)

	os.Setenv("WEATHER", "955781466c1e414e9e9181300240806")
	tempClient, err := api.NewWeatherClient(&http.Client{}, os.Getenv("WEATHER"))
	if err != nil {
		t.Fatal(err)
		return
	}

	mySpan := new(traceSpan)
	mySpan.On("End").Return()

	myTracerMock := &tracerMock{
		span: mySpan,
	}
	myTracerMock.On("Start", mock.Anything, mock.Anything, mock.Anything).Times(3).Return()

	tempUsecase := usecase.NewTempUsecase(tempClient)
	kelvinService := usecase.NewKelvinService()
	tempByCEPctrl := internal.NewTempByLocaleController(myTracerMock, localidadeUsecase, tempUsecase, kelvinService)

	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	req, err := http.NewRequest("GET", "/cep/91210290", nil)
	if err != nil {
		t.Fatal(err)
	}

	chiCtx := chi.NewRouteContext()

	// Create a new test request with the additional Chi contetx
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))

	// Add the key/value to the context.
	chiCtx.URLParams.Add("cep", fmt.Sprintf("%v", "91210290"))

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := MakeCepHandler(tempByCEPctrl)
	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)
	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected := `invalid zipcode`
	if rr.Body.String() == expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}
