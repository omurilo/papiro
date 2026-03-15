package main

import (
	"log"

	"github.com/omurilo/papiro/internal/builder"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new Papiro blog",
	Run: func(cmd *cobra.Command, args []string) {
		if err := builder.InitSite(); err != nil {
			log.Fatalf("erro ao inicializar o site: %v", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
