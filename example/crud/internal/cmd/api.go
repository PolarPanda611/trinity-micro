// Author: Daniel TAN
// Date: 2021-10-02 01:20:48
// LastEditors: Daniel TAN
// LastEditTime: 2021-10-04 01:18:10
// FilePath: /trinity-micro/example/crud/cmd/api.go
// Description:
package cmd

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/PolarPanda611/trinity-micro"
	"github.com/PolarPanda611/trinity-micro/core/dbx"
	"github.com/PolarPanda611/trinity-micro/core/logx"
	"github.com/PolarPanda611/trinity-micro/core/tracerx"
	_ "github.com/PolarPanda611/trinity-micro/example/crud/internal/adapter/controller"
	"github.com/PolarPanda611/trinity-micro/example/crud/internal/application/model"
	"github.com/PolarPanda611/trinity-micro/example/crud/internal/config"
	"github.com/PolarPanda611/trinity-micro/example/crud/internal/consts"
	"gorm.io/gorm"

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

// @title           trinity-micro Example API
// @version         1.0
// @description     This is a sample server for trinity-micro
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:3000
// @BasePath  /

// @securityDefinitions.basic  BasicAuth
func RunAPI(cmd *cobra.Command, args []string) {
	serviceName := fmt.Sprintf("%v-%v", consts.ProjectName, consts.ApiCmdString)

	// infra set up
	logger := logx.Init(logx.Config{
		ServiceName: serviceName,
		LogfilePath: fmt.Sprintf("%v.log", serviceName),
	})
	log.Printf("%v-%v service starting!\n", consts.ProjectName, consts.ApiCmdString)
	logger.Infof("%v-%v service starting!", consts.ProjectName, consts.ApiCmdString)

	currentPath, _ := os.Getwd()
	configPath := filepath.Join(currentPath + "/conf/config.toml")
	if _, err := toml.DecodeFile(configPath, &config.Conf); err != nil {
		logger.Fatalf("load config :%v failed, err: %v", configPath, err)
	}
	logger.Infof("load config: %v successfully", config.Conf)

	dbx.Init(&dbx.Config{
		Type:        config.Conf.Database.Type,
		DSN:         config.Conf.Database.DSN,
		TablePrefix: config.Conf.Database.TablePrefix,
		MaxIdleConn: config.Conf.Database.MaxIdleConn,
		MaxOpenConn: config.Conf.Database.MaxOpenConn,
		Logger:      logger.WithField("app", "database"),
	})
	// handle multi tenant initialize
	{
		tenants := make([]string, 0)
		sessionDB := dbx.DB.Session(&gorm.Session{
			NewDB: true,
		})
		sessionDB.AutoMigrate(&model.Tenant{})
		var res []model.Tenant
		if err := sessionDB.Find(&res).Error; err != nil {
			logger.Fatalf("list tenant failed, err: ", err)
		}
		for _, tenant := range res {
			tenants = append(tenants, fmt.Sprintf("tn_%d", tenant.ID))
		}
		dbx.Migrate(context.Background(), tenants...)
	}

	tracerx.Init(tracerx.Config{
		Type:        config.Conf.Tracer.Type,
		ServiceName: config.Conf.Tracer.ServiceName,
		Host:        config.Conf.Tracer.Host,
	})
	t := trinity.New(trinity.Config{
		Logger: logger,
	})
	if err := t.Start(":3000"); err != nil {
		log.Printf("%v-%v service terminated, error:%v \n", consts.ProjectName, consts.ApiCmdString, err)
		logger.Fatalf("%v-%v service terminated, error:%v", consts.ProjectName, consts.ApiCmdString, err)
	}

}
