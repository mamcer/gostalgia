package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	scanCmd = &cobra.Command{
		Use:   "scan",
		Short: "scans a specific directory and commits info to nostalgia database",
		Long:  ``,
		Run:   scan,
	}
)

func scan(ccmd *cobra.Command, args []string) {
	fmt.Printf("hello there nostalgia config: %v, tags: %v\n", viper.Get("nostalgia_path"), tags)
}
