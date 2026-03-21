package main

import (
	"log"

	"github.com/omurilo/papiro/internal/builder"
	"github.com/spf13/cobra"
)

var path string

var initCmd = &cobra.Command{
	Use:   "init [path]",
	Short: "Initialize a new Papiro blog at path (default \".\")",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 0 {
			path = args[0]
		}
		if err := builder.InitSite(path); err != nil {
			log.Fatalf("erro ao inicializar o site: %v", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
