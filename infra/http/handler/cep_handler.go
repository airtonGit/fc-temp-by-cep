package handler

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"regexp"

	"github.com/airtongit/fc-temp-by-cep/internal"
	"github.com/airtongit/fc-temp-by-cep/internal/usecase"
	"github.com/go-chi/chi/v5"
)

type TempByCep interface {
	GetTemp(ctx context.Context, cep string) (internal.Temp, error)
}

func validate(cep string) error {
	log.Println("validate CEP", cep)
	matched, err := regexp.MatchString(`^\d{8}$`, cep)
	if err != nil {
		return err
	}
	if !matched {
		log.Println("matchstring not match with", cep)
		return errors.New("invalid zipcode")
	}
	return nil
}

func MakeCepHandler(ctrl TempByCep) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("handling req")
		cep := chi.URLParam(r, "cep")
		if err := validate(cep); err != nil {
			w.WriteHeader(http.StatusUnprocessableEntity)
			w.Write([]byte(err.Error()))
			return
		}
		log.Println("cep", cep)
		tempResponse, err := ctrl.GetTemp(r.Context(), cep)
		if err != nil {
			log.Println("get_temp err", err)
			if errors.Is(err, usecase.ErrCepNotFound) {
				w.WriteHeader(http.StatusNotFound)
				w.Write([]byte(err.Error()))
				return
			}

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
