package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"regexp"

	"github.com/airtongit/fc-temp-by-cep/internal"
	"github.com/go-chi/chi"
)

type TempByCep interface {
	GetTemp(ctx context.Context, cep string) (internal.Temp, error)
}

func validate(cep string) error {
	matched, err := regexp.MatchString(`ˆ\d{8}`, cep)
	if err != nil {
		return err
	}
	if !matched {
		return errors.New("invalid zipcode")
	}
	return nil
}

func CepHandler(ctrl TempByCep) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cep := chi.URLParam(r, "cep")
		if err := validate(cep); err != nil {
			w.WriteHeader(http.StatusUnprocessableEntity)
			w.Write([]byte(err.Error()))
			return
		}

		tempResponse, err := ctrl.GetTemp(r.Context(), cep)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		w.Header().Set("Content-Type", "application/json")

		err = json.NewEncoder(w).Encode(tempResponse)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
	}
}
