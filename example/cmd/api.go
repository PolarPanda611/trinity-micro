package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

var (
	apiCmd = &cobra.Command{
		Use:   apiCmdString,
		Short: fmt.Sprintf("starting the %v service for %v", apiCmdString, projectName),
		Long:  fmt.Sprintf("This is the %v service for %v", apiCmdString, projectName),
		Run: func(cmd *cobra.Command, args []string) {
			log.Printf("%v:%v service starting ", projectName, apiCmdString)
		},
	}
)

func init() {
	rootCmd.AddCommand(apiCmd)
}
