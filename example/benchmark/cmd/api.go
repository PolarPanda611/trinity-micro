package cmd

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/PolarPanda611/trinity-micro"
	"github.com/PolarPanda611/trinity-micro/example/benchmark/internal/adapter/controller"
	_ "github.com/PolarPanda611/trinity-micro/example/benchmark/internal/adapter/controller"
	"github.com/PolarPanda611/trinity-micro/example/benchmark/internal/consts"

	"github.com/PolarPanda611/trinity-micro/example/benchmark/internal/infra/containers"

	"github.com/PolarPanda611/trinity-micro/example/benchmark/internal/infra/logx"

	"github.com/go-chi/chi/v5"
	"github.com/spf13/cobra"
)

var (
	apiCmd = &cobra.Command{
		Use:   consts.ApiCmdString,
		Short: fmt.Sprintf("starting the %v service for %v", consts.ApiCmdString, consts.ProjectName),
		Long:  fmt.Sprintf("This is the %v service for %v", consts.ApiCmdString, consts.ProjectName),
		Run:   RunAPI,
	}
)

func init() {
	rootCmd.AddCommand(apiCmd)
}

// @title Trinity Micro Example API
// @version 1.0
// @description This is a sample server celler server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:3000
// @BasePath /api/v1
// @query.collection.format multi

// @securityDefinitions.basic BasicAuth

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

// @securitydefinitions.oauth2.application OAuth2Application
// @tokenUrl https://example.com/oauth/token
// @scope.write Grants write access
// @scope.admin Grants read and write access to administrative information

// @securitydefinitions.oauth2.implicit OAuth2Implicit
// @authorizationurl https://example.com/oauth/authorize
// @scope.write Grants write access
// @scope.admin Grants read and write access to administrative information

// @securitydefinitions.oauth2.password OAuth2Password
// @tokenUrl https://example.com/oauth/token
// @scope.read Grants read access
// @scope.write Grants write access
// @scope.admin Grants read and write access to administrative information

// @securitydefinitions.oauth2.accessCode OAuth2AccessCode
// @tokenUrl https://example.com/oauth/token
// @authorizationurl https://example.com/oauth/authorize
// @scope.admin Grants read and write access to administrative information

// @x-extension-openapi {"example": "value on a json format"}
func RunAPI(cmd *cobra.Command, args []string) {
	log.Printf("%v:%v service starting ", consts.ProjectName, consts.ApiCmdString)
	// infra set up
	logx.Init()
	containers.Init()

	// init Router
	r := chi.NewRouter()
	trinity.DIRouter(r, containers.Container)
	r.Get("/benchmark/simple_raw", controller.SimpleRaw)
	logx.Logger.Infof("request mapping: method: %-6s %-30s => handler: %v ", "GET", "/benchmark/simple_raw", "SimpleRaw")
	r.Get("/benchmark/path_param_raw/{id}", controller.PathParamRaw)
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
	if err := s.ListenAndServe(); err != nil {
		logx.Logger.Fatalf("service terminated, error:%v", err)
	}
}
