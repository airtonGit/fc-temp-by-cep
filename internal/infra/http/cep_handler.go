package http

import (
	"fmt"
	"github.com/go-chi/chi"
	"net/http"
)

func CepHandler(w http.ResponseWriter, r *http.Request) {
	cep := chi.URLParam(r, "cep")
	w.Write([]byte(fmt.Sprintf("Hello World! %s", cep)))
}
