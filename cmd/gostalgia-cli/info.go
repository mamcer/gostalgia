package main

import (
	"fmt"
	"github.com/spf13/cobra"
)

var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "Show current configuration and paths",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Nostalgia Scan Path: '%s'\n", cfg.NOSTALGIA_SCAN_PATH)
		fmt.Printf("Nostalgia Home Path: '%s'\n", cfg.NOSTALGIA_HOME_PATH)
		fmt.Printf("Nostalgia Thumb Target Path: '%s'\n", cfg.NOSTALGIA_THUMB_TARGET_PATH)
		fmt.Printf("Nostalgia DB Connection: '%s'\n", cfg.NOSTALGIA_CONNECTION_STRING)
	},
}

func init() {
	rootCmd.AddCommand(infoCmd)
}
