// Author: Daniel TAN
// Date: 2021-08-18 00:07:41
// LastEditors: Daniel TAN
// LastEditTime: 2021-10-02 00:47:02
// FilePath: /trinity-micro/example/cmd/root.go
// Description:
package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/PolarPanda611/trinity-micro/example/internal/consts"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   consts.ProjectName,
	Short: fmt.Sprintf("%v command line tool", consts.ProjectName),
	Long:  fmt.Sprintf("%v command line tool, generated by trinity-micro ", consts.ProjectName),
	Run: func(cmd *cobra.Command, args []string) {
		log.Printf("%v version: %v", consts.ProjectName, consts.Version)
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("root cmd execute failed, error:%v", err)
		os.Exit(1)
	}
}
