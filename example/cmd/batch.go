// Author: Daniel TAN
// Date: 2021-08-18 00:34:11
// LastEditors: Daniel TAN
// LastEditTime: 2021-10-02 00:47:24
// FilePath: /trinity-micro/example/cmd/batch.go
// Description:
/*
 * @Author: Daniel TAN
 * @Date: 2021-08-18 00:34:11
 * @LastEditors: Daniel TAN
 * @LastEditTime: 2021-09-09 00:08:47
 * @FilePath: /trinity-micro/example/cmd/batch.go
 * @Description:
 */
package cmd

import (
	"fmt"
	"log"

	"github.com/PolarPanda611/trinity-micro/example/internal/consts"
	"github.com/spf13/cobra"
)

var (
	batchCmd = &cobra.Command{
		Use:   consts.BatchCmdString,
		Short: fmt.Sprintf("starting the %v service for %v", consts.BatchCmdString, consts.ProjectName),
		Long:  fmt.Sprintf("This is the %v service for %v", consts.BatchCmdString, consts.ProjectName),
		Run:   RunBatch,
	}
)

func init() {
	rootCmd.AddCommand(batchCmd)
}

func RunBatch(cmd *cobra.Command, args []string) {
	log.Printf("%v:%v service starting ", consts.ProjectName, consts.BatchCmdString)
}
