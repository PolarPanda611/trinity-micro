// Author: Daniel TAN
// Date: 2021-10-02 00:49:54
// LastEditors: Daniel TAN
// LastEditTime: 2021-10-02 01:23:42
// FilePath: /trinity-micro/example/benchmark/internal/adapter/controller/raw.go
// Description:
package controller

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/PolarPanda611/trinity-micro/core/httpx"
	"github.com/go-chi/chi/v5"
)

func SimpleRaw(w http.ResponseWriter, r *http.Request) {
	res := httpx.SuccessResponse{
		Status: 200,
		Result: "ok",
	}
	b, _ := json.Marshal(res)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(b)
}
func PathParamRaw(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, _ := strconv.Atoi(idStr)
	res := httpx.SuccessResponse{
		Status: 200,
		Result: id,
	}
	b, _ := json.Marshal(res)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(b)
}
