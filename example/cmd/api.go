package cmd

import (
	"fmt"
	"log"
	"net/http"
	"time"
	"trinity-micro/core/ioc/container"
	_ "trinity-micro/example/internal/adapter/controller"
	"trinity-micro/example/internal/infra/containers"

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
	log.Printf("%v:%v service starting ", projectName, apiCmdString)
	containers.Container.InstanceDISelfCheck()
	r := chi.NewRouter()
	container.DIRouter(r, containers.Container)
	// r.Get("/test/user2", func(w http.ResponseWriter, r *http.Request) {
	// 	res := httpx.Response{
	// 		Status: 200,
	// 		Result: "haha",
	// 	}
	// 	b, _ := json.Marshal(res)
	// 	w.WriteHeader(200)
	// 	w.Write(b)
	// })
	// r.Get("/test/user2/{id}", func(w http.ResponseWriter, r *http.Request) {
	// 	idStr := chi.URLParam(r, "id")
	// 	id, _ := strconv.Atoi(idStr)
	// 	res := httpx.Response{
	// 		Status: 200,
	// 		Result: id,
	// 	}
	// 	b, _ := json.Marshal(res)
	// 	w.WriteHeader(200)
	// 	w.Write(b)
	// })
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
