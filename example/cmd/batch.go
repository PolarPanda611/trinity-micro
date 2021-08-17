package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

var (
	batchCmd = &cobra.Command{
		Use:   batchCmdString,
		Short: fmt.Sprintf("starting the %v service for %v", batchCmdString, projectName),
		Long:  fmt.Sprintf("This is the %v service for %v", batchCmdString, projectName),
		Run: func(cmd *cobra.Command, args []string) {
			log.Printf("%v:%v service starting ", projectName, batchCmdString)
		},
	}
)

func init() {
	rootCmd.AddCommand(batchCmd)
}
