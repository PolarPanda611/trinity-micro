package cmd

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/spf13/cobra"
)

var (
	apiCmd = &cobra.Command{
		Use:   apiCmdString,
		Short: fmt.Sprintf("starting the %v service for %v", apiCmdString, projectName),
		Long:  fmt.Sprintf("This is the %v service for %v", apiCmdString, projectName),
		Run: func(cmd *cobra.Command, args []string) {
			log.Printf("%v:%v service starting ", projectName, apiCmdString)
			r := chi.NewRouter()
			r.Use()
			r.Get("/", func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte("hi"))
			})
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
		},
	}
)

func init() {
	rootCmd.AddCommand(apiCmd)
}
