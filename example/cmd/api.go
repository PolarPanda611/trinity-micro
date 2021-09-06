package cmd

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
	"trinity-micro/core/httpx"
	"trinity-micro/core/ioc/container"
	_ "trinity-micro/example/internal/adapter/controller"
	"trinity-micro/example/internal/infra/containers"
	"trinity-micro/example/internal/infra/logx"

	"github.com/go-chi/chi/v5"
	"github.com/spf13/cobra"
)

var (
	apiCmd = &cobra.Command{
		Use:   apiCmdString,
		Short: fmt.Sprintf("starting the %v service for %v", apiCmdString, projectName),
		Long:  fmt.Sprintf("This is the %v service for %v", apiCmdString, projectName),
		Run:   RunAPI,
	}
)

func init() {
	rootCmd.AddCommand(apiCmd)
}

func RunAPI(cmd *cobra.Command, args []string) {
	logx.Logger.Infof("%v:%v service starting ", projectName, apiCmdString)
	r := chi.NewRouter()
	container.DIRouter(r, containers.Container)
	r.Get("/benchmark/simple_raw", SimpleRaw)
	logx.Logger.Infof("request mapping: method: %-6s %-30s => handler: %v ", "GET", "/benchmark/simple_raw", "SimpleRaw")
	r.Get("/benchmark/path_param_raw/{id}", PathParamRaw)
	logx.Logger.Infof("request mapping: method: %-6s %-30s => handler: %v ", "GET", "/benchmark/path_param_raw/{id}", "SimpleRaw")
	s := &http.Server{
		Addr:              ":3000",
		Handler:           r,
		ReadTimeout:       time.Duration(10) * time.Second,
		ReadHeaderTimeout: time.Duration(10) * time.Second,
		WriteTimeout:      time.Duration(10) * time.Second,
		IdleTimeout:       time.Duration(10) * time.Second,
		MaxHeaderBytes:    5 * 1024 * 1024,
	}
	s.ListenAndServe()
}

func SimpleRaw(w http.ResponseWriter, r *http.Request) {
	res := httpx.Response{
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
	res := httpx.Response{
		Status: 200,
		Result: id,
	}
	b, _ := json.Marshal(res)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(b)
}
