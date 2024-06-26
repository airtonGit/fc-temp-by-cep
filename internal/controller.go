package internal

import (
	"context"
	"fmt"
	"github.com/airtongit/fc-temp-by-cep/internal/usecase"
	"go.opentelemetry.io/otel/trace"
	"log"
)

type LocalidadeUsecase interface {
	Execute(ctx context.Context, input usecase.LocalidadeInput) (usecase.LocalidadeOutput, error)
}

type TempUsecase interface {
	Execute(ctx context.Context, input usecase.TempUsecaseInput) (usecase.TempUsecaseOutput, error)
}

type KelvinService interface {
	GetKelvin(tempC float64) float64
}

type TraceSpan interface {
	End(options ...trace.SpanEndOption)
}

type TraceAdapter interface {
	Start(ctx context.Context, spanName string, opts ...trace.SpanStartOption) (context.Context, TraceSpan)
}

type spanAdapter struct {
	otelSpan trace.Span
}

func (s *spanAdapter) End(options ...trace.SpanEndOption) {
	s.otelSpan.End(options...)
}

func NewSpanAdapter(span trace.Span) *spanAdapter {
	return &spanAdapter{otelSpan: span}
}

type tracerAdapter struct {
	otelTracer trace.Tracer
}

func NewTracerAdapter(tracer trace.Tracer) *tracerAdapter {
	return &tracerAdapter{otelTracer: tracer}
}

func (t *tracerAdapter) Start(ctx context.Context, spanName string, opts ...trace.SpanStartOption) (context.Context, TraceSpan) {
	ctx, span := t.otelTracer.Start(ctx, spanName, opts...)
	return ctx, NewSpanAdapter(span)
}

type tempByLocaleController struct {
	localidadeUsecase LocalidadeUsecase
	tempUsecase       TempUsecase
	kelvinService     KelvinService
	OTELTracer        TraceAdapter
}

func NewTempByLocaleController(tracer TraceAdapter, localidadeUsecase LocalidadeUsecase, tempUsecase TempUsecase, kelvinService KelvinService) *tempByLocaleController {
	return &tempByLocaleController{
		localidadeUsecase: localidadeUsecase,
		tempUsecase:       tempUsecase,
		kelvinService:     kelvinService,
		OTELTracer:        tracer,
	}
}

type Temp struct {
	TempC      float64 `json:"temp_C,omitempty"`
	TempF      float64 `json:"temp_F,omitempty"`
	TempK      float64 `json:"temp_K,omitempty"`
	Localidade string  `json:"localidade"`
}

func (t *tempByLocaleController) GetTemp(ctx context.Context, cep string) (Temp, error) {

	ctx, spanInicial := t.OTELTracer.Start(ctx, "get_temp_by_cep full")
	defer spanInicial.End()

	ctx, spanLocalidade := t.OTELTracer.Start(ctx, "Chama externa get localidade by CEP")
	defer spanLocalidade.End()

	localidadeInput := usecase.LocalidadeInput{
		Cep: cep,
	}

	log.Println("localidade_usecase_exec", localidadeInput)

	localidade, err := t.localidadeUsecase.Execute(ctx, localidadeInput)
	if err != nil {
		if err.Error() == usecase.ErrCepNotFound.Error() {
			log.Println("ctrl error ir err_cep_not_found")
			return Temp{}, usecase.ErrCepNotFound
		}
		return Temp{}, fmt.Errorf("getting localidade by cep: %w", err)
	}

	tempUsecaseInput := usecase.TempUsecaseInput{
		Localidade: localidade.Localidade,
		Uf:         localidade.Uf,
		Pais:       localidade.Pais,
	}

	spanLocalidade.End()

	ctx, spanTemperatura := t.OTELTracer.Start(ctx, "Chama externa: temperatura by CEP")
	defer spanTemperatura.End()

	log.Println("temperature_usecase_exec", localidadeInput)
	temp, err := t.tempUsecase.Execute(ctx, tempUsecaseInput)
	if err != nil {
		return Temp{}, fmt.Errorf("getting temp by localidade: %w", err)
	}

	kelvin := t.kelvinService.GetKelvin(temp.TempC)

	return Temp{
		TempC:      temp.TempC,
		TempF:      temp.TempF,
		TempK:      kelvin,
		Localidade: fmt.Sprintf("%s %s", localidade.Localidade, localidade.Uf),
	}, nil
}
